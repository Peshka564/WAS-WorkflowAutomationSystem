package main

import (
	"database/sql"
	"fmt"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/models"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/repositories"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    // parseTime = true -> parses DATETIME into time.Time
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true")
    if err != nil {
        fmt.Println("Could not connect to db");
        fmt.Println(err);
        return;
    }

    repo := repositories.Workflow{Db: db}
    newWorkflow := models.Workflow{Name: "Pesho", Active: true}
    repo.Insert(&newWorkflow)
    fmt.Println(newWorkflow.Id)
    // r := chi.NewRouter()
    // r.Use(middleware.Logger)
    // r.Use(cors.Handler(cors.Options{
    //     AllowedOrigins: []string{"http://*"},
    //     AllowedMethods: []string{"GET"},
    // }))
    // r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, Mama!"))
    // })
    // http.ListenAndServe(":3000", r)
}