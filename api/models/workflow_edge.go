package models

type WorkflowEdge struct {
	BaseModel
	NodeFrom   int
	NodeTo     int
	WorkflowId int
}