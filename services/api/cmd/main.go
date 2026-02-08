package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/server"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/services"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
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

    userConn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to User Service: %v", err)
	}
	defer userConn.Close()
	userService := services.User{ GrpcClient: pb.NewUserServiceClient(userConn) }

    app := server.App{
        Db: db,
        Router: chi.NewRouter(),
        Validator: validator.New(),
        UserService: &userService,
    }

    app.Router.Use(middleware.Logger)
    app.Router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://*"},
        AllowedMethods: []string{"GET", "POST"},
    }))
    
    app.Router.Post("/api/workflows/create", app.CreateWorkflow)
    app.Router.Post("/api/register", app.RegisterUser)
    app.Router.Post("/api/login", app.LoginUser)
    
    err = http.ListenAndServe(":3000", app.Router)
    if err != nil {
        log.Fatal("Could not start server", err);
        return;
    }
}