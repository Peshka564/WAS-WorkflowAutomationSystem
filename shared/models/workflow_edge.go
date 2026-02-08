package models

import "time"

type WorkflowEdge struct {
	Id         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	
	NodeFrom   string
	NodeTo     string
	WorkflowId int
}