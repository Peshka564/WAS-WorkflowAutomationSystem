package services

import (
	"context"
	"encoding/json"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/api/repositories"
	api_utils "github.com/Peshka564/WAS-WorkflowAutomationSystem/api/utils"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"
	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
)

type Workflow struct {
	WorkflowRepo repositories.Workflow
	WorkflowNodeRepo repositories.WorkflowNode
	WorkflowEdgeRepo repositories.WorkflowEdge
}

func (service *Workflow) CreateWorkflow(ctx context.Context, data dto.CreateWorkflowPayload) error {
	// TODO: Unit of Work/Transaction
	newWorkflow, err := service.insertWorkflow(ctx, data.Workflow)
	if err != nil {
		return err
	}
	newNodes, err := service.insertNodes(ctx, data.Nodes, newWorkflow)
	if err != nil {
		return err
	}
	return service.insertEdges(ctx, data, newNodes, newWorkflow)
}

func (service *Workflow) insertWorkflow(ctx context.Context, workflow dto.CreateWorkflow) (*models.Workflow, error) {
	newWorkflow := &models.Workflow{Name: workflow.Name, Active: false, UserId: 0}
	err := service.WorkflowRepo.Insert(newWorkflow)
	if err != nil {
		return nil, err
	}
	return newWorkflow, nil
}

func (service *Workflow) insertNodes(ctx context.Context, nodes []dto.CreateWorkflowNode, newWorkflow *models.Workflow) ([]models.WorkflowNode, error) {
	transformedNodes := make([]models.WorkflowNode, 0)
	for _, node := range nodes {
		var transformedNode models.WorkflowNode
		transformedNode.TaskName = node.TaskName
		transformedNode.WorkflowId = newWorkflow.Id
		transformedNode.Type = models.Listener
		jsonPos, _ := json.Marshal(node.Position)
		transformedNode.Position = string(jsonPos)
		
		err := service.WorkflowNodeRepo.Insert(&transformedNode)
		if err != nil {
			return nil, err
		}
		transformedNodes = append(transformedNodes, transformedNode)
	}
	return transformedNodes, nil
}

func (service *Workflow) insertEdges(ctx context.Context, data dto.CreateWorkflowPayload, nodes []models.WorkflowNode, workflow *models.Workflow) error {
	newEdges := make([]models.WorkflowEdge, 0)
	for _,edge := range data.Edges {
		nodeFrom, idxFrom := api_utils.Find(data.Nodes, func(node dto.CreateWorkflowNode) bool { return node.DisplayId == edge.From})
		if nodeFrom == nil {
			return errs.InvalidInputError{}
		}
		nodeTo, idxTo := api_utils.Find(data.Nodes, func(node dto.CreateWorkflowNode) bool { return node.DisplayId == edge.To})
		if nodeTo == nil {
			return errs.InvalidInputError{}
		}
		nodeFromId := nodes[idxFrom].Id
		nodeToId := nodes[idxTo].Id
		newEdges = append(newEdges, models.WorkflowEdge{NodeFrom: nodeFromId, NodeTo: nodeToId, WorkflowId: workflow.Id})
	}
	return service.WorkflowEdgeRepo.InsertMany(newEdges)
}