package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userService "github.com/Peshka564/WAS-WorkflowAutomationSystem/services/user"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	DB *sql.DB
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// TODO: Repository and a whole separate user service

	res, err := s.DB.ExecContext(ctx, "INSERT INTO users (username, name, password_hash) VALUES (?, ?, ?)", req.Username, req.Name, hashedPassword)
	if err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            return nil, status.Error(codes.AlreadyExists, "username already exists")
        }
        return nil, err
    }

	userId, err := res.LastInsertId()
	if err != nil {
		// TODO: Check
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}

	token, _ := userService.GenerateJWT(userId, os.Getenv("JWT_SECRET"))

	return &pb.AuthResponse{
		Token: token,
		User:  &pb.User{Id: userId, Username: req.Username, Name: req.Name},
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	var user models.User
	err := s.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE username = ?", req.Username).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Name, &user.PasswordHash)
	
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, _ := userService.GenerateJWT(int64(user.Id), os.Getenv("JWT_SECRET"))

	return &pb.AuthResponse{
		Token: token,
		User:  &pb.User{Id: int64(user.Id), Username: req.Username, Name: user.Name},
	}, nil
}

func (s *UserServiceServer) GetCredentials(ctx context.Context, req *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	var credential models.Credential

	err := s.DB.QueryRowContext(ctx, `SELECT * FROM credentials WHERE id = ?`, req.CredentialId).Scan(&credential.Id, &credential.ServiceName, &credential.UserId, &credential.AccessToken, &credential.RefreshToken, &credential.ExpiresAt)

	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.NotFound, "connection not found")
	}

	// If not expired for more than 5 mins
	if time.Now().Add(5 * time.Minute).Before(credential.ExpiresAt) {
		return &pb.GetCredentialsResponse{AccessToken: credential.AccessToken, Success: true}, nil
	}

	// If expired, try to refresh

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
	}

	token := &oauth2.Token{
		RefreshToken: credential.RefreshToken,
		Expiry:       time.Now().Add(-1 * time.Hour), // Force expiry to trigger refresh
	}

	// Refresh
	tokenSource := conf.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()

	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		// TODO: Add active field to the credentials
		return nil, status.Error(codes.Unauthenticated, "connection revoked, please reconnect")
	}

	// Check for refresh token rotation
	if newToken.RefreshToken != "" {
		credential.RefreshToken = newToken.RefreshToken
	}

	_, err = s.DB.ExecContext(ctx, `
		UPDATE credentials 
		SET access_token = ?, refresh_token = ?, expires_at = ? 
		WHERE id = ?`,
		newToken.AccessToken, credential.RefreshToken, newToken.Expiry, req.CredentialId)
	if err != nil {
		log.Printf("Failed to save new token: %v", err)
	}

	return &pb.GetCredentialsResponse{AccessToken: newToken.AccessToken, Success: true}, nil
}

func main() {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true")
    if err != nil {
        log.Fatal("Could not connect to db", err);
        return;
    }

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Could not load ENV vars", err);
		return;
	}

	listener, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServiceServer{DB: db})

	log.Printf("User Service running on :50055...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}