package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"
)

type Workflow struct {
	GrpcClient pb.WorkflowServiceClient
}

func (s *Workflow) CreateWorkflow(ctx context.Context, data dto.CreateWorkflowPayload) (*dto.CreateWorkflowResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Println(data.Nodes)
	nodeInput := make([]*pb.NodeInput, 0)
	for _, node := range data.Nodes {
		var nodeId *string
		if node.Id != nil {
			nodeId = node.Id
		}

		nodeInput = append(nodeInput, &pb.NodeInput{
			Id:          nodeId,
			DisplayId:   node.DisplayId,
			ServiceName: node.ServiceName,
			TaskName:    node.TaskName,
			Type:        node.Type,
			Position:    node.Position,
			Config:      node.Config,
			CredentialId: node.CredentialId,
		})
	}

	edgeInput := make([]*pb.EdgeInput, 0)
	for _, edge := range data.Edges {
		var edgeId *string
		if edge.Id != nil {
			edgeId = edge.Id
		}
		edgeInput = append(edgeInput, &pb.EdgeInput{
			Id:        edgeId,
			DisplayId: edge.DisplayId,
			FromId:    edge.From,
			ToId:      edge.To,
		})
	}

	rawId := ctx.Value("user_id")
    if rawId == nil {
        return nil, errors.New("internal error: user_id missing from context")
    }

	var workflowId int64 = 0
	if data.Workflow.Id != nil {
		workflowId = int64(*data.Workflow.Id)
	}

	req := &pb.CreateWorkflowRequest{
		Id:     workflowId, // NEW: 0 = Create, >0 = Update
		Name:   data.Workflow.Name,
		UserId: rawId.(int64),
		Nodes:  nodeInput,
		Edges:  edgeInput,
	}

	res, err := s.GrpcClient.CreateWorkflow(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.CreateWorkflowResponse{ WorkflowId: int(res.Id) }, nil
}

func (s *Workflow) GetWorkflows(ctx context.Context) ([]dto.Workflow, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    userId, ok := ctx.Value("user_id").(int64)
    if !ok {
        return nil, fmt.Errorf("user_id not found in context")
    }

    req := &pb.GetWorkflowsRequest{
        UserId: userId,
    }

    res, err := s.GrpcClient.GetWorkflows(ctx, req)
    if err != nil {
        return nil, err
    }

    workflows := make([]dto.Workflow, 0, len(res.Workflows))
    
    for _, w := range res.Workflows {
        workflows = append(workflows, dto.Workflow{
            Id:        int(w.Id),
            Name:      w.Name,
            Active:    w.IsActive,
            CreatedAt: w.CreatedAt.AsTime(),
            UpdatedAt: w.UpdatedAt.AsTime(),
            UserId:    int(w.UserId),
        })
    }

    return workflows, nil
}

func (s *Workflow) GetWorkflowById(ctx context.Context, workflowId int) (*dto.GetWorkflowResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req := &pb.GetWorkflowByIdRequest{
        Id: int64(workflowId),
    }

    res, err := s.GrpcClient.GetWorkflowById(ctx, req)
    if err != nil {
        return nil, err
    }

	fmt.Println(res)

	workflow := dto.Workflow{
		Id: int(res.Workflow.Id),
		CreatedAt: res.Workflow.CreatedAt.AsTime(),
		UpdatedAt: res.Workflow.UpdatedAt.AsTime(),
		Name: res.Workflow.Name,
		Active: res.Workflow.IsActive,
		UserId: int(res.Workflow.UserId),
	}

	nodes := make([]dto.GetNodeResponse, 0)
	for _, node := range res.Nodes {
		nodes = append(nodes, dto.GetNodeResponse{
			Id: node.Id,
			DisplayId:  node.DisplayId,
			ServiceName: node.ServiceName ,
			TaskName:     node.TaskName,
			WorkflowId:   int(node.WorkflowId),
			Type:         node.Type,
			Position:     node.Position,
			Config:       node.Config,
			CredentialId: node.CredentialId,
		})
	}

	edges := make([]dto.GetEdgeResponse, 0)
	for _, edge := range res.Edges {
		edges = append(edges, dto.GetEdgeResponse{
			Id: edge.Id,
			DisplayId:  edge.DisplayId,
			WorkflowId:  int(edge.WorkflowId),
			NodeFrom: edge.FromId,
			NodeTo: edge.ToId,
		})
	}

    return &dto.GetWorkflowResponse{
		Workflow: workflow,
		Nodes: nodes,
		Edges: edges,
	}, nil
}
