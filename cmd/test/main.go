package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

func main() {
	// TODO: Validate when using service discovery
	orchConn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer orchConn.Close()

	s := pb.NewOrchestratorClient(orchConn)
	res, err := s.TriggerWorkflow(context.Background(), &pb.TriggerRequest{
		WorkflowId: 1,
		InitialPayload: `{ "hello": 1}`,
	})
	fmt.Println(err)
	fmt.Println(res)
}