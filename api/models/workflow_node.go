package models

type WorkflowNodeType int

const (
	Listener WorkflowNodeType = iota
	Action
	Transformer
)

func (nt WorkflowNodeType) String() string {
	switch nt {
	case Listener:
		return "listener"
	case Action:
		return "action"
	case Transformer:
		return "transformer"
	default:
		panic("Invalid workflow node type")
	}
}

type WorkflowNode struct {
	BaseModel
	WorkflowId   int
	TaskName     string
	WorkflowType WorkflowNodeType
	Position     string // JSON encoded position { x: ..., y: ... }
}