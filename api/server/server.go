package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/dto"
	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/errors"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/repositories"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type App struct {
	Router  *chi.Mux
	Db *sql.DB
	Validator *validator.Validate
}

func (app *App) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateWorkflowPayload
	err := parseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		sendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	fmt.Println("OPAAAAAAAAAAAAAAAAAAAAAAAAAa")
	err = app.Validator.Struct(payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, formValidationErrorMessage(err))
		return
	}
	fmt.Println("HEREEEEEEEEEEEEEEEEEEEEEEEEEEEEEEee")

	workflowService := services.Workflow{
		WorkflowRepo: repositories.Workflow{Db: app.Db},
		WorkflowNodeRepo: repositories.WorkflowNode{Db: app.Db},
		WorkflowEdgeRepo: repositories.WorkflowEdge{Db: app.Db},
	}
	err = workflowService.CreateWorkflow(r.Context(), payload)
	if err != nil {
		if errors.Is(err, errs.InvalidInputError{}) {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}