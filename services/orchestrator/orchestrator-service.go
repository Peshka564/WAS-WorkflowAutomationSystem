package orchestrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/api/repositories"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/utils"
)

type OrchestratorService struct {
	Db           *sql.DB
	GmailService pb.TaskWorkerClient
	UserService pb.UserServiceClient
	// service -> grpc address
	// Registry     map[string]string
}

type ExecutionContext struct {
	WorkflowID   int
	// The "Bag of State"
	// TODO: Store this in DB
	CurrentData  map[string]interface{}
}

func (orchestrator *OrchestratorService) ExecuteWorkflow(ctx context.Context, listenerNodeId string, initialPayload string) error {
	workflowNodeRepo := repositories.WorkflowNode{ Db: orchestrator.Db }

	listenerNode, err := workflowNodeRepo.FindById(listenerNodeId)
	if err != nil {
		return fmt.Errorf("Invalid trigger node");
	}
	workflowId := listenerNode.WorkflowId

	log.Printf("Starting Workflow %d", workflowId)

	var triggerData map[string]interface{}
	if err := json.Unmarshal([]byte(initialPayload), &triggerData); err != nil {
		return fmt.Errorf("failed to parse initial payload: %v", err)
	}

	state := &ExecutionContext{
		WorkflowID:  workflowId,
		CurrentData: map[string]interface{}{"trigger": triggerData},
	}

	nodes, err := orchestrator.getNodesInLinearOrder(listenerNode)
	if err != nil {
		return err
	}
	

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
			return execErr 
		}

		// 4. Update State with Results
		// So Step 2 can access {{ step_1.data }}
		var outputMap map[string]interface{}
		if outputJSON != "" {
			if err := json.Unmarshal([]byte(outputJSON), &outputMap); err != nil {
				return fmt.Errorf("failed to parse outputJSON: %v", err)
			}
			state.CurrentData[node.Id] = outputMap
		}
	}

	log.Printf("Workflow %d Completed Successfully", workflowId)
	return nil
}

func (orchestrator *OrchestratorService) executeAction(ctx context.Context, node models.WorkflowNode, state *ExecutionContext) (string, error) {
	// A. Variable Resolution (The "Smart Orchestrator" Pattern)
	// We take '{"subject": "Hello {{trigger.name}}"}' and turn it into '{"subject": "Hello Bob"}'
	// resolvedConfig, err := e.resolveVariables(node.ConfigJSON, state.CurrentData)
	// if err != nil {
	// 	return "", fmt.Errorf("variable resolution failed: %v", err)
	// }

	// Authenticate with the third party api for the task
	var authToken string
	if node.CredentialId != nil {
		tokenResp, err := orchestrator.UserService.GetCredentials(ctx, &pb.GetCredentialsRequest{
			CredentialId: int32(*node.CredentialId),
		})
		if err != nil {
			return "", fmt.Errorf("auth failure: %v", err)
		}
		authToken = tokenResp.AccessToken
	}

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
			ConfigJson: `{
				"to": "petaryp@uni-sofia.bg",
				"subject": "Answer",
				"body": "Answer body"
			}`,
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

func (orchestrator *OrchestratorService) getNodesInLinearOrder(listenerNode *models.WorkflowNode) ([]models.WorkflowNode, error) {
	// TODO: Get from workflow service
	nodeRepo := repositories.WorkflowNode{ Db: orchestrator.Db }
	nodes, err := nodeRepo.FindByWorkflowId(listenerNode.WorkflowId)
	if err != nil {
        return nil, fmt.Errorf("failed to fetch nodes: %v", err)
    }

	idToNode := make(map[string]models.WorkflowNode)
	for _, node := range nodes {
		idToNode[node.Id] = node
    }

	adjList := make(map[string][]string)
	
	edgeRepo := repositories.WorkflowEdge{ Db: orchestrator.Db }
	edges, err := edgeRepo.FindByWorkflowId(listenerNode.WorkflowId)
	if err != nil {
        return nil, fmt.Errorf("failed to fetch edges: %v", err)
    }

	for _, edge := range edges {
		adjList[edge.NodeFrom] = append(adjList[edge.NodeFrom], edge.NodeTo)
	}

	// TODO: Only toposort the part of the graph that is reachable from listenerNode
	toposorted, err := utils.TopologicalSort(adjList)
	if err != nil {
		return nil, err
	}

	// Remove the trigger
	toposorted = toposorted[1:]

	toposortedNodes := make([]models.WorkflowNode, 0)
	for _, nodeId := range toposorted {
		toposortedNodes = append(toposortedNodes, idToNode[nodeId])
	}

	return toposortedNodes, nil
}

// func (orchestrator *OrchestratorService) getWorkflowStatus(workflowId int): active {

// }