package models

import "time"

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
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	WorkflowId   int
	ServiceName  string
	ActionName   string
	Type         WorkflowNodeType
	Config       string // JSON encoded
	CredentialId int
	Position     string // JSON encoded position { x: ..., y: ... }
}