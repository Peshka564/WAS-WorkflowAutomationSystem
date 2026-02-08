package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WorkflowServiceServer struct {
	pb.UnimplementedWorkflowServiceServer
	Db *sql.DB
}

func (s *WorkflowServiceServer) GetWorkflows(ctx context.Context, req *pb.GetWorkflowsRequest) (*pb.GetWorkflowsResponse, error) {
	workflowRepo := repositories.Workflow{ Db: s.Db }

	dbWorkflows, err := workflowRepo.FindByUserId(req.UserId)
	if err != nil {
		log.Printf("Repo error: %v", err)
		return nil, status.Error(codes.Internal, "failed to fetch workflows")
	}

	var workflows []*pb.Workflow
	for _, w := range dbWorkflows {
		workflows = append(workflows, &pb.Workflow{
			Id:        int64(w.Id),
			Name:      w.Name,
			IsActive:  w.Active,
			UserId:    int64(w.UserId),
			CreatedAt: timestamppb.New(w.CreatedAt),
			UpdatedAt: timestamppb.New(w.UpdatedAt),
		})
	}

	return &pb.GetWorkflowsResponse{Workflows: workflows}, nil
}

func (s *WorkflowServiceServer) CreateWorkflow(ctx context.Context, req *pb.CreateWorkflowRequest) (*pb.CreateWorkflowResponse, error) {
	workflowRepo := repositories.Workflow{ Db: s.Db }
	workflowNodeRepo := repositories.WorkflowNode{ Db: s.Db }
	workflowEdgeRepo := repositories.WorkflowEdge{ Db: s.Db }

	// TODO: Transaction + Unit of Work
	
	workflowModel := &models.Workflow{
		Name:   req.Name,
		UserId: int(req.UserId),
		Active: false,
	}

	if err := workflowRepo.Insert(workflowModel); err != nil {
		log.Printf("Failed to insert workflow: %v", err)
		return nil, status.Error(codes.Internal, "failed to create workflow")
	}

	workflowId := workflowModel.Id
	nodeIdMap := make(map[string]string)
	for _, nodeReq := range req.Nodes {
		nodeModel := &models.WorkflowNode{
			Id:         uuid.New().String(),
			WorkflowId: workflowId,
			ServiceName: nodeReq.ServiceName,
			TaskName:   nodeReq.TaskName,
			Type:       models.FromString(nodeReq.Type),
			Config: 	nodeReq.Config,
			CredentialId: nodeReq.CredentialId,
			DisplayId: nodeReq.DisplayId,
			Position:   nodeReq.Position,
		}

		fmt.Println(nodeModel)
		if err := workflowNodeRepo.Insert(nodeModel); err != nil {
			log.Printf("Failed to insert node %s: %v", nodeReq.DisplayId, err)
			return nil, status.Error(codes.Internal, "failed to save workflow nodes")
		}
		nodeIdMap[nodeReq.DisplayId] = nodeModel.Id
	}

	var edgesToInsert []models.WorkflowEdge
	for _, edgeReq := range req.Edges {
		fromId, okFrom := nodeIdMap[edgeReq.FromId]
		toId, okTo := nodeIdMap[edgeReq.ToId]

		if !okFrom || !okTo {
			return nil, status.Errorf(codes.InvalidArgument, "edge references unknown node: %s -> %s", edgeReq.FromId, edgeReq.ToId)
		}
		edgesToInsert = append(edgesToInsert, models.WorkflowEdge{
			Id:         uuid.New().String(),
			WorkflowId: workflowId,
			DisplayId:  edgeReq.DisplayId,
			NodeFrom:   fromId,
			NodeTo:     toId,
		})
	}

	if len(edgesToInsert) > 0 {
		if err := workflowEdgeRepo.InsertMany(edgesToInsert); err != nil {
			log.Printf("Failed to insert edges: %v", err)
			return nil, status.Error(codes.Internal, "failed to save workflow connections")
		}
	}

	return &pb.CreateWorkflowResponse{Id: int64(workflowId)}, nil
}

func (s *WorkflowServiceServer) ActivateWorkflow(ctx context.Context, req *pb.ActivateWorkflowRequest) (*pb.ActivateWorkflowResponse, error) {
	workflowRepo := repositories.Workflow{ Db: s.Db }
    err := workflowRepo.UpdateActiveStatus(int(req.Id), req.Active)
    if err != nil {
        return &pb.ActivateWorkflowResponse{Success: false}, err
    }

    return &pb.ActivateWorkflowResponse{Success: true}, nil
}

func (s *WorkflowServiceServer) GetWorkflowById(ctx context.Context, req *pb.GetWorkflowByIdRequest) (*pb.GetWorkflowByIdResponse, error) {
	workflowRepo := repositories.Workflow{ Db: s.Db }
	workflowNodeRepo := repositories.WorkflowNode{ Db: s.Db }
	workflowEdgeRepo := repositories.WorkflowEdge{ Db: s.Db }

    workflow, err := workflowRepo.FindById(int(req.Id))
	if err != nil {
		// TODO: Return status error from google
		return nil, err
	}
	fmt.Println(workflow)
    nodes, err := workflowNodeRepo.FindByWorkflowId(int(req.Id))
	if err != nil {
		return nil, err
	}
	fmt.Println(nodes)
    edges, err := workflowEdgeRepo.FindByWorkflowId(int(req.Id))
	if err != nil {
		return nil, err
	}
	fmt.Println(edges)

	nodesMapped := make([]*pb.Node, 0)
	for _, node := range nodes {
		nodesMapped = append(nodesMapped, &pb.Node{
			Id: node.Id,
			DisplayId: node.DisplayId,
			Config: node.Config,
			WorkflowId: int64(node.WorkflowId),
			TaskName: node.TaskName,
			ServiceName: node.ServiceName,
			Type: node.Type.String(),
			Position: node.Position,
			CredentialId: node.CredentialId,
		})
	}

	edgesMapped := make([]*pb.Edge, 0)
	for _, edge := range edges {
		edgesMapped = append(edgesMapped, &pb.Edge{
			Id: edge.Id,
			DisplayId: edge.DisplayId,
			WorkflowId: int64(edge.WorkflowId),
			FromId: edge.NodeFrom,
			ToId: edge.NodeTo,
		})
	}

    return &pb.GetWorkflowByIdResponse{ Workflow: &pb.Workflow{
		Name: workflow.Name,
		Id: int64(workflow.Id),
		UserId: int64(workflow.UserId),
		IsActive: workflow.Active,
		CreatedAt: timestamppb.New(workflow.CreatedAt),
		UpdatedAt: timestamppb.New(workflow.UpdatedAt),
	}, Nodes: nodesMapped, Edges: edgesMapped,}, nil
}


func main() {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true")
    if err != nil {
        log.Fatal("Could not connect to db", err);
        return;
    }

	listener, err := net.Listen("tcp", ":50056")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWorkflowServiceServer(grpcServer, &WorkflowServiceServer{Db: db})

	log.Printf("Workflow Service running on :50056...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}