package models

type Workflow struct {
	BaseModel
	Name   string
	Active bool
	UserId int
}