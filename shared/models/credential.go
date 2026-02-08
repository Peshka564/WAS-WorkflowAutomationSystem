package models

import "time"

type Credential struct {
	Id        		int
	ServiceName 	string
	UserId 			int
	AccessToken 	string
	RefreshToken 	string
	ExpiresAt 		time.Time
}