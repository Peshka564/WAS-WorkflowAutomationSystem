package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/middleware"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/server"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/services"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    // parseTime = true -> parses DATETIME into time.Time
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

    userConn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to User Service: %v", err)
	}
	defer userConn.Close()
	userService := services.User{ GrpcClient: pb.NewUserServiceClient(userConn) }

    workflowConn, err := grpc.NewClient("localhost:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to Workflow Service: %v", err)
	}
	defer workflowConn.Close()
	workflowService := services.Workflow{ GrpcClient: pb.NewWorkflowServiceClient(workflowConn) }


	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/api/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send", "https://www.googleapis.com/auth/gmail.readonly"},
		Endpoint:     google.Endpoint,
	}

    app := server.App{
        Db: db,
        Router: chi.NewRouter(),
        Validator: validator.New(),
        UserService: &userService,
        WorkflowService: &workflowService,
		OAuthConfig: googleOauthConfig,
    }

    app.Router.Use(chi_middleware.Logger)
    app.Router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://*"},
        AllowedMethods: []string{"GET", "POST", "OPTIONS", "PATCH"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    }))
    
    app.Router.Post("/api/register", app.RegisterUser)
    app.Router.Post("/api/login", app.LoginUser)
	app.Router.Get("/api/auth/google/callback", app.GoogleCallback)
    
    app.Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/api/workflows", app.GetWorkflows)
		r.Post("/api/workflows", app.CreateWorkflow)
        r.Patch("/api/workflows/{id}/activate", app.ActivateWorkflow)
        r.Get("/api/workflows/{id}", app.GetWorkflowById)
		r.Get("/api/connections", app.GetConnections)
		r.Get("/api/auth/google/login", app.GoogleLogin)
		r.Get("/api/templates", app.GetTemplates)
		r.Post("/api/templates", app.SaveTemplate)
	})
    
    err = http.ListenAndServe(":3000", app.Router)
    if err != nil {
        log.Fatal("Could not start server", err);
        return;
    }
}