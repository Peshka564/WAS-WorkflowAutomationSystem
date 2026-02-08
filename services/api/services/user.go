package services

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"
	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

type User struct {
	GrpcClient pb.UserServiceClient
}

func (s *User) Register(ctx context.Context, user dto.RegisterUserPayload) (*dto.UserWithTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.RegisterRequest{
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
	}

	data, err := s.GrpcClient.Register(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return nil, errs.AlreadyExists{EntityName: "User"}
		}
		return nil, err
	}

	res := dto.UserWithTokenResponse{
		Token: data.Token,
		User: dto.UserResponse{
			Id: int(data.User.Id),
			Username: data.User.Username,
			Name: data.User.Name,
		},
	}

	return &res, nil
}

func (s *User) Login(ctx context.Context, user dto.LoginUserPayload) (*dto.UserWithTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	}

	data, err := s.GrpcClient.Login(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			return nil, errs.NotFoundError{EntityName: "User"}
		}
		return nil, err
	}

	res := dto.UserWithTokenResponse{
		Token: data.Token,
		User: dto.UserResponse{
			Id: int(data.User.Id),
			Username: data.User.Username,
			Name: data.User.Name,
		},
	}

	return &res, nil
}