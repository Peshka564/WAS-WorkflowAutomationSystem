package models

import "time"

type Workflow struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	
	Name      string
	Active    bool
	UserId    int
}