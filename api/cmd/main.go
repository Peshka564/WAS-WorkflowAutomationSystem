package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    // parseTime = true -> parses DATETIME into time.Time
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true")
    if err != nil {
        log.Fatal("Could not connect to db", err);
        return;
    }

    app := server.App{
        Db: db,
        Router: chi.NewRouter(),
        Validator: validator.New(),
    }
    app.Router.Use(middleware.Logger)
    app.Router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://*"},
        AllowedMethods: []string{"GET", "POST"},
    }))
    app.Router.Post("/workflows/create", app.CreateWorkflow)
    
    err = http.ListenAndServe(":3000", app.Router)
    if err != nil {
        log.Fatal("Could not start server", err);
        return;
    }
}