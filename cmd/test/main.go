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
	userConn, _ := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer userConn.Close()

	s := pb.NewUserServiceClient(userConn)
	res, err := s.Login(context.Background(), &pb.LoginRequest{
		Username: "pesho",
		Password: "pesho",
	})
	fmt.Println(err)
	fmt.Println(res)
}