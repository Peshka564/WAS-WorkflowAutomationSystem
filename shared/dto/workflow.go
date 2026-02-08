package dto

import "time"

type CreateWorkflow struct {
	Name string
}

type CreateWorkflowNode struct {
	Id 			 *string `json:"id"`
	DisplayId    string `validate:"required"`
	ServiceName  string `validate:"required"`
	TaskName     string `validate:"required"`
	Type         string `validate:"required,oneof=listener action transformer"`
	Position     string `validate:"required"`
	Config       string `validate:"required"`
	CredentialId *int32 `json:"credential_id"`
}

type CreateWorkflowEdge struct {
	From      string `validate:"required"`
	To        string `validate:"required"`
	DisplayId string `validate:"required"`
}

type CreateWorkflowPayload struct {
	Workflow CreateWorkflow       `validate:"required"`
	Nodes    []CreateWorkflowNode `validate:"required,min=1,dive"`
	Edges    []CreateWorkflowEdge `validate:"required,dive"`
}

type CreateWorkflowResponse struct {
	WorkflowId int `json:"workflowId"`
}

type Workflow struct {
	Id        int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string `json:"name"`
	Active bool `json:"active"`
	UserId int `json:"user_id"`
}

type ActivateWorkflowPayload struct {
    Active bool `json:"active"`
}

type GetNodeResponse struct {
	Id 			 string `json:"id"`
	DisplayId    string `json:"display_id"`
	ServiceName  string `json:"service_name"`
	TaskName     string `json:"task_name"`
	WorkflowId   int 	`json:"workflow_id"`
	Type         string `json:"type"`
	Position     string `json:"position"`
	Config       string `json:"config"`
	CredentialId *int32 `json:"credential_id"`
}

type GetEdgeResponse struct {
	Id 			 string `json:"id"`
	DisplayId    string `json:"display_id"`
	WorkflowId   int 	`json:"workflow_id"`
	NodeFrom 	 string `json:"node_from"`
	NodeTo 	 string `json:"node_to"`
}

type GetWorkflowResponse struct {
	Workflow Workflow          `json:"workflow"`
	Nodes    []GetNodeResponse `json:"nodes"`
	Edges    []GetEdgeResponse `json:"edges"`
}
