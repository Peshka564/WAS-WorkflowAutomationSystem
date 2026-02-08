package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"

	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"golang.org/x/oauth2"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type GmailServer struct {
	pb.UnimplementedTaskWorkerServer
}

type EmailConfig struct {
	To      string
	Subject string
	Body    string
}

func (s *GmailServer) ExecuteTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	if(req.TaskName == "send_email") {
		var config EmailConfig
		// TODO: Validator
		if err := json.Unmarshal([]byte(req.ConfigJson), &config); err != nil {
			return &pb.TaskResponse{Success: false, ErrorMessage: "Invalid config JSON"}, nil
		}

		if req.AuthToken == "" {
			return &pb.TaskResponse{Success: false, ErrorMessage: "Missing OAuth2 Access Token"}, nil
		}

		// TODO: Validate
		client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
            &oauth2.Token{AccessToken: req.AuthToken},
        ))

		srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return &pb.TaskResponse{Success: false, ErrorMessage: fmt.Sprintf("Gmail Client Error: %v", err)}, nil
		}

		messageStr := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s", config.To, config.Subject, config.Body)

		gMessage := &gmail.Message{
			Raw: base64.URLEncoding.EncodeToString([]byte(messageStr)),
		}

		_, err = srv.Users.Messages.Send("me", gMessage).Do()
		if err != nil {
			// TODO: Handle Google API specific errors (e.g., 401 Unauthorized, 403 Quota Exceeded)
			return &pb.TaskResponse{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Google API Error: %v", err),
			}, nil
		}
		return &pb.TaskResponse{Success: true, OutputPayload: `{"status": "sent"}`}, nil
	}
	return &pb.TaskResponse{Success: false, ErrorMessage: "Invalid task name"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskWorkerServer(s, &GmailServer{})

	log.Println("Gmail Service running on :50052")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}