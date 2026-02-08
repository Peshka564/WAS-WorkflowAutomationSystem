package dto

type CreateWorkflow struct {
	Name string
}

type NodePosition struct {
	X float32 `validate:"required"`
	Y float32 `validate:"required"`
}

type CreateWorkflowNode struct {
	DisplayId     string       `validate:"required"`
	TaskName      string       `validate:"required"`
	Position      NodePosition `validate:"required"`
	ThirdPartyApp string       `json:"service" validate:"oneof=gmail drive github"`
}

type CreateWorkflowEdge struct {
	From string `validate:"required"`
	To   string `validate:"required"`
}

type CreateWorkflowPayload struct {
	Workflow CreateWorkflow       `validate:"required"`
	Nodes    []CreateWorkflowNode `validate:"required,min=1,dive"`
	Edges    []CreateWorkflowEdge `validate:"required,dive"`
}