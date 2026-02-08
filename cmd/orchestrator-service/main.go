package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/orchestrator"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

// This is the gRPC adapter around the orchestrator
type OrchestratorServiceServer struct {
	pb.UnimplementedOrchestratorServer
	OrchestratorService *orchestrator.OrchestratorService
}

func (s *OrchestratorServiceServer) TriggerWorkflow(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerResponse, error) {
	// Note: In a real system, use RabbitMQ or some other message broker here
	go func() {
		err := s.OrchestratorService.ExecuteWorkflow(context.Background(), int(req.WorkflowId), req.InitialPayload)
		if err != nil {
			log.Printf("Background execution failed: %v", err)
		}
	}()

	return &pb.TriggerResponse{
		ExecutionId: 1, // TODO: Get from DB, maybe uuid
		Success:     true,
	}, nil
}

func main() {
	// parseTime = true -> parses DATETIME into time.Time
	// TODO: Change this to some other port
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true")
    if err != nil {
        log.Fatal("Could not connect to db", err);
        return;
    }

	// TODO: Use Service Discovery/Registry

	// TODO: Validate when using service discovery
	gmailConn, _ := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer gmailConn.Close()
	// userConn, _ := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// userClient := pb.NewUserServiceClient(userConn)

	orchestrator := &orchestrator.OrchestratorService{
		DB: db,
		GmailService: pb.NewTaskWorkerClient(gmailConn),
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, &OrchestratorServiceServer{OrchestratorService: orchestrator})

	log.Println("Orchestrator Service running on :50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}