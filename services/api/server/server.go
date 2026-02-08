package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/repositories"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/services"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"

	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type App struct {
	Router  *chi.Mux
	Db *sql.DB
	Validator *validator.Validate
	UserService *services.User
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
	err = app.Validator.Struct(payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, formValidationErrorMessage(err))
		return
	}

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

func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserPayload
	err := parseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		sendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, formValidationErrorMessage(err))
		return
	}

	res, err := app.UserService.Register(r.Context(), payload)
	if err != nil {
		if(errors.Is(err, errs.AlreadyExists{EntityName: "User"})) {
			sendError(w, http.StatusUnauthorized, "User already exists");
			return
		}
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginUserPayload
	err := parseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		sendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, formValidationErrorMessage(err))
		return
	}

	res, err := app.UserService.Login(r.Context(), payload)
	if err != nil {
		if(errors.Is(err, errs.NotFoundError{EntityName: "User"})) {
			sendError(w, http.StatusUnauthorized, "User not found");
			return
		}
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}