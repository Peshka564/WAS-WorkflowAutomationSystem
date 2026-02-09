package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/services"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/utils"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/dto"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
	"golang.org/x/oauth2"

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
	OAuthConfig *oauth2.Config
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

func generateState(userID int64) string {
    data := fmt.Sprintf("%d", userID)
    h := hmac.New(sha256.New, []byte(os.Getenv("OAUTH_STATE_SECRET")))
    h.Write([]byte(data))
    signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
    return fmt.Sprintf("%s|%s", data, signature)
}

func verifyState(state string) (int64, error) {
    parts := strings.Split(state, "|")
    if len(parts) != 2 {
        return 0, errors.New("invalid state format")
    }
    userIDStr, signature := parts[0], parts[1]

    // Verify Signature
    h := hmac.New(sha256.New, []byte(os.Getenv("OAUTH_STATE_SECRET")))
    h.Write([]byte(userIDStr))
    expectedSig := base64.URLEncoding.EncodeToString(h.Sum(nil))

    if signature != expectedSig {
        return 0, errors.New("invalid state signature")
    }

    var userID int64
    fmt.Sscanf(userIDStr, "%d", &userID)
    return userID, nil
}

func (app *App) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
    if !ok {
        utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
	
	state := generateState(userID)

	url := app.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "url": url,
    })
}

func (app *App) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
    if state == "" {
        utils.SendError(w, http.StatusBadRequest, "State parameter missing")
        return
    }

	userID, err := verifyState(state)
    if err != nil {
        fmt.Println("State Verification Failed:", err)
        utils.SendError(w, http.StatusUnauthorized, "Invalid OAuth state")
        return
    }

	code := r.URL.Query().Get("code")
	
	token, err := app.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Code exchange failed")
		return
	}

	err = app.saveCredential(userID, "gmail", token)
	if err != nil {
        fmt.Println(err)
		utils.SendError(w, http.StatusInternalServerError, "Failed to save credentials")
		return
	}

	http.Redirect(w, r, "http://localhost:5173/connections?status=success", http.StatusSeeOther)
}

func (app *App) saveCredential(userID int64, service string, token *oauth2.Token) error {
	// Upsert logic (Update if exists, else Insert)
	query := `
		INSERT INTO credentials (user_id, service_name, access_token, refresh_token, expires_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			access_token = VALUES(access_token),
			refresh_token = VALUES(refresh_token),
			expires_at = VALUES(expires_at)
	`
	_, err := app.Db.Exec(query, userID, service, token.AccessToken, token.RefreshToken, token.Expiry)
	return err
}

func (app *App) GetConnections(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(int64)
    
    rows, err := app.Db.Query("SELECT service_name FROM credentials WHERE user_id = ?", userID)
    if err != nil {
		fmt.Println(err)
        utils.SendError(w, http.StatusInternalServerError, "Db Error")
        return
    }
    defer rows.Close()

    connections := []map[string]interface{}{}
    for rows.Next() {
        var service string
        rows.Scan(&service)
        connections = append(connections, map[string]interface{}{
            "service": service,
            "connected": true,
        })
    }
    
    json.NewEncoder(w).Encode(connections)
}

func (app *App) GetTemplates(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	rows, err := app.Db.Query("SELECT id, name, subject, body, email_to FROM email_templates WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	templates := []dto.Template{}
	for rows.Next() {
		var t dto.Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Subject, &t.Body, &t.EmailTo); err != nil {
			continue
		}
		templates = append(templates, t)
	}

	json.NewEncoder(w).Encode(templates)
}

func (app *App) SaveTemplate(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(int64)
    var t dto.Template
    
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        utils.SendError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := app.Validator.Struct(t); err != nil {
        utils.SendError(w, http.StatusBadRequest, "Validation failed")
        return
    }

    if t.ID == 0 {
        query := "INSERT INTO email_templates (user_id, name, subject, body, email_to) VALUES (?, ?, ?, ?, ?)"
        res, err := app.Db.Exec(query, userID, t.Name, t.Subject, t.Body, t.EmailTo)
        if err != nil {
            utils.SendError(w, http.StatusInternalServerError, "Failed to create template")
            return
        }

        newID, _ := res.LastInsertId()
        t.ID = int(newID)
        w.WriteHeader(http.StatusCreated)

    } else {
        query := "UPDATE email_templates SET name=?, subject=?, body=?, email_to=? WHERE id=? AND user_id=?"
        res, err := app.Db.Exec(query, t.Name, t.Subject, t.Body, t.EmailTo, t.ID, userID)
        if err != nil {
            utils.SendError(w, http.StatusInternalServerError, "Failed to update template")
            return
        }

        rowsAffected, _ := res.RowsAffected()
        if rowsAffected == 0 {
            utils.SendError(w, http.StatusNotFound, "Template not found or unauthorized")
            return
        }
        w.WriteHeader(http.StatusOK)
    }

    json.NewEncoder(w).Encode(t)
}