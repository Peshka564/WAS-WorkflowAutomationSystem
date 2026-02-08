package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/services"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/utils"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"

	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type App struct {
	Router  *chi.Mux
	Db *sql.DB
	Validator *validator.Validate
	UserService *services.User
	WorkflowService *services.Workflow
}

func (app *App) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	res, err := app.WorkflowService.GetWorkflows(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	fmt.Println(res)
	w.WriteHeader(http.StatusOK)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}

func (app *App) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateWorkflowPayload
	err := utils.ParseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.FormValidationErrorMessage(err))
		return
	}

	res, err := app.WorkflowService.CreateWorkflow(r.Context(), payload)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}

func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.FormValidationErrorMessage(err))
		return
	}

	res, err := app.UserService.Register(r.Context(), payload)
	if err != nil {
		if(errors.Is(err, errs.AlreadyExists{EntityName: "User"})) {
			utils.SendError(w, http.StatusUnauthorized, "User already exists");
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginUserPayload
	err := utils.ParseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.FormValidationErrorMessage(err))
		return
	}

	res, err := app.UserService.Login(r.Context(), payload)
	if err != nil {
		if(errors.Is(err, errs.NotFoundError{EntityName: "User"})) {
			utils.SendError(w, http.StatusUnauthorized, "User not found");
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}

func (app *App) ActivateWorkflow(w http.ResponseWriter, r *http.Request) {
    workflowIDStr := chi.URLParam(r, "id")
    workflowID, err := strconv.Atoi(workflowIDStr)
    if err != nil {
        http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
        return
    }

    _, ok := r.Context().Value("user_id").(int64)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	var payload dto.ActivateWorkflowPayload
	err = utils.ParseJSON(r, &payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	err = app.Validator.Struct(payload)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.FormValidationErrorMessage(err))
		return
	}

    req := &pb.ActivateWorkflowRequest{
        Id:     int64(workflowID),
        Active: payload.Active,
    }

    resp, err := app.WorkflowService.GrpcClient.ActivateWorkflow(r.Context(), req)
    if err != nil || !resp.Success {
		fmt.Println(err)
        http.Error(w, "Failed to activate workflow", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (app *App) GetWorkflowById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEEEEEEEEEEEEEEEEEEEST")
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.SendError(w, http.StatusBadRequest, "Invalid workflow ID")
        return
    }

	// userId := r.Context().Value("user_id").(int64)

	fmt.Println("TEST")

    res, err := app.WorkflowService.GetWorkflowById(r.Context(), id)
    if err != nil {
		fmt.Println(err)
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
    }

    w.WriteHeader(http.StatusOK)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Write(jsonRes)
}


// func (app *App) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
//     idStr := chi.URLParam(r, "id")
//     id, err := strconv.Atoi(idStr)
//     if err != nil {
//         utils.SendError(w, http.StatusBadRequest, "Invalid workflow ID")
//         return
//     }

//     userID, ok := r.Context().Value("user_id").(int64)
//     if !ok {
//         utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
//         return
//     }

//     var payload dto.CreateWorkflowPayload
//     if err := utils.ParseJSON(r, &payload); err != nil {
//         utils.SendError(w, http.StatusBadRequest, "Invalid JSON")
//         return
//     }

//     if err := app.Validator.Struct(payload); err != nil {
//         utils.SendError(w, http.StatusBadRequest, utils.FormValidationErrorMessage(err))
//         return
//     }

//     var nodeInput []*pb.NodeInput
//     for _, node := range payload.Nodes {
//         nodeInput = append(nodeInput, &pb.NodeInput{
//             DisplayId:   node.DisplayId,
//             ServiceName: node.ServiceName,
//             TaskName:    node.TaskName,
//             Type:        node.Type,
//             Position:    node.Position,
//             Config:      node.Config,
//             CredentialId: node.CredentialId,
//         })
//     }

//     var edgeInput []*pb.Edge
//     for _, e := range payload.Edges {
//         edgeInput = append(edgeInput, &pb.Edge{
//             From:      e.From,
//             To:        e.To,
//             DisplayId: e.DisplayId,
//         })
//     }

//     // 6. Call gRPC Service
//     _, err = app.WorkflowService.GrpcClient.UpdateWorkflow(r.Context(), &pb.UpdateWorkflowRequest{
//         Id:     int64(id),
//         UserId: userID,
//         Name:   payload.Workflow.Name,
//         Nodes:  nodeInput,
//         Edges:  edgeInput,
//     })

//     if err != nil {
//         fmt.Println("Update Error:", err)
//         utils.SendError(w, http.StatusInternalServerError, "Failed to update workflow")
//         return
//     }

//     w.WriteHeader(http.StatusOK)
//     json.NewEncoder(w).Encode(map[string]bool{"success": true})
// }