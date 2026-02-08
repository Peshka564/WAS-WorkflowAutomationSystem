package orchestrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

type OrchestratorService struct {
	DB           *sql.DB
	GmailService pb.TaskWorkerClient
	// service -> grpc address
	// Registry     map[string]string
}

// ExecutionContext holds the state of a running workflow
type ExecutionContext struct {
	WorkflowID   int
	CurrentData  map[string]interface{} // The "Bag of State" (trigger.body, step_1.data, etc)
}

func (orchestrator *OrchestratorService) ExecuteWorkflow(ctx context.Context, workflowID int, initialPayload string) error {
	log.Printf("Starting Workflow %d", workflowID)

	// 1. Initialize State
	// We parse the initial trigger data (e.g., Webhook body) into the map
	var triggerData map[string]interface{}
	if err := json.Unmarshal([]byte(initialPayload), &triggerData); err != nil {
		return fmt.Errorf("failed to parse initial payload: %v", err)
	}

	state := &ExecutionContext{
		WorkflowID:  workflowID,
		CurrentData: map[string]interface{}{"trigger": triggerData},
	}

	// 2. Load the Workflow Graph (Simplified: Fetching all nodes)
	// In a real app, you would fetch edges to determine the order. 
	// Here we assume nodes have a 'next_node_id' or we fetch them in order.
	nodes, err := orchestrator.getNodesInLinearOrder(workflowID)
	if err != nil {
		return err
	}

	// 3. Execution Loop
	for _, node := range nodes {
		log.Printf("Executing Node: %s (%s)", node.Id, node.Type)

		var outputJSON string
		var execErr error

		switch node.Type.String() {
		case "action":
			outputJSON, execErr = orchestrator.executeAction(ctx, node, state)
		// case "transformer":
		// 	outputJSON, execErr = orchestrator.executeTransformer(ctx, node, state)
		default:
			log.Printf("Skipping unknown node type: %s", node.Type)
			continue
		}

		// Handle Failures
		if execErr != nil {
			log.Printf("Workflow Failed at Node %s: %v", node.Id, execErr)
			// TODO: Add Retry Logic or 'Dead Letter Queue' here
			return execErr 
		}

		// 4. Update State with Results
		// So Step 2 can access {{ step_1.data }}
		var outputMap map[string]interface{}
		if outputJSON != "" {
			json.Unmarshal([]byte(outputJSON), &outputMap)
			state.CurrentData[node.Id] = outputMap
		}
	}

	log.Printf("Workflow %d Completed Successfully", workflowID)
	return nil
}

func (orchestrator *OrchestratorService) executeAction(ctx context.Context, node models.WorkflowNode, state *ExecutionContext) (string, error) {
	// A. Variable Resolution (The "Smart Orchestrator" Pattern)
	// We take '{"subject": "Hello {{trigger.name}}"}' and turn it into '{"subject": "Hello Bob"}'
	// resolvedConfig, err := e.resolveVariables(node.ConfigJSON, state.CurrentData)
	// if err != nil {
	// 	return "", fmt.Errorf("variable resolution failed: %v", err)
	// }

	// B. Authentication (Zapier Style)
	// If this node requires a connection (Gmail, Slack), fetch the fresh token
	// var authToken string
	// if node.ConnectionID != nil {
	// 	// Call User Service to get fresh token (handles refresh logic transparently)
	// 	tokenResp, err := e.UserService.GetConnectionToken(ctx, &pb.GetTokenRequest{
	// 		ConnectionId: int32(*node.ConnectionID),
	// 	})
	// 	if err != nil {
	// 		return "", fmt.Errorf("auth failure: %v", err)
	// 	}
	// 	authToken = tokenResp.AccessToken
	// }
	// C. Service Discovery
	// Look up where the worker lives (e.g., "gmail" -> "localhost:50052")
	// serviceAddr, exists := e.Registry[node.ServiceSlug]
	// if !exists {
	// 	// Fallback: Query DB if not in cache
	// 	serviceAddr = "localhost:50052" // Hardcoded fallback for demo
	// }

	if(node.ServiceName == "gmail") {
		res, err := orchestrator.GmailService.ExecuteTask(ctx, &pb.TaskRequest{
			TaskName:   node.TaskName,
			ConfigJson: node.Config,
			AuthToken:  authToken,
		})
		if err != nil {
			return "", err
		}
		if !res.Success {
			return "", fmt.Errorf("Task failed: %s", res.ErrorMessage)
		}

		return res.OutputPayload, nil
	}
	return "", nil
	// D. gRPC Call
	// Dial the worker dynamically
	// conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return "", err
	// }
	// defer conn.Close()

	// client := pb.NewTaskWorkerClient(conn)

	// resp, err := client.ExecuteTask(ctx, &pb.TaskRequest{
	// 	TaskType:   node.ActionSlug,   // e.g., "send_email"
	// 	ConfigJson: resolvedConfig,    // The CLEAN, resolved JSON
	// 	AuthToken:  authToken,         // The OAuth Access Token
	// })

	// if err != nil {
	// 	return "", err
	// }

	// if !resp.Success {
	// 	return "", fmt.Errorf("remote task failed: %s", resp.ErrorMessage)
	// }

	// return resp.OutputPayload, nil
}

// ---------------------------------------------------------------------------
// 3. Transformer Executor (The "Internal Logic")
// ---------------------------------------------------------------------------

// func (e *Engine) executeTransformer(ctx context.Context, node Node, state *ExecutionContext) (string, error) {
// 	// Example: A "Delay" transformer
// 	if node.ActionSlug == "delay" {
// 		// Parse config to find duration
// 		// e.g. {"seconds": 5}
// 		log.Println("Transformer: Sleeping for 2 seconds...")
// 		time.Sleep(2 * time.Second)
// 		return `{"status": "delayed"}`, nil
// 	}
	
// 	return "{}", nil
// }

// ---------------------------------------------------------------------------
// 4. Helpers (Parser & DB)
// ---------------------------------------------------------------------------

// func (e *Engine) resolveVariables(templateJSON string, data map[string]interface{}) (string, error) {
// 	// Uses Go's text/template to replace {{ key }} with values
// 	tmpl, err := template.New("node").Option("missingkey=error").Parse(templateJSON)
// 	if err != nil {
// 		return "", err
// 	}

// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

func (orchestrator *OrchestratorService) getNodesInLinearOrder(workflowID int) ([]models.WorkflowNode, error) {
	// TODO: SELECT * FROM workflow_nodes WHERE workflow_id = ?
	
	return []models.WorkflowNode{
		{
			Id:           "1",
			WorkflowId: workflowID,
			ServiceName:  "gmail",
			TaskName:   "send_email",
			Type:         models.WorkflowNodeType(1),
			Config:   `{"to": "petaryp@uni-sofia.bg", "subject": "Testing orchestrator", "body": "Testing orchestrator body"}`,
			CredentialId: 1,
		},
	}, nil
}

// func (orchestrator *OrchestratorService) getWorkflowStatus(workflowId int): active {

// }