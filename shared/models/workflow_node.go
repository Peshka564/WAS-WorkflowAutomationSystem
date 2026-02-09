package models

import (
	"time"
)

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


func FromString(s string) WorkflowNodeType {
	switch s {
	case "listener":
		return Listener
	case "action":
		return Action
	case "transformer":
		return Transformer
	default:
		panic("Invalid workflow node type")
	}
}

type WorkflowNode struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	DisplayId string
	WorkflowId   int
	ServiceName  string
	TaskName   string
	Type         WorkflowNodeType
	Config       string // JSON encoded
	CredentialId *int32
	Position     string // JSON encoded position { x: ..., y: ... }
}