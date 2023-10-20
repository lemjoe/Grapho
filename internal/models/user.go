package models

import "time"

type User struct {
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Password string `json:"passwd"`

	Email        string    `json:"email"`
	IsAdmin      bool      `json:"is_admin"`
	CreationDate time.Time `json:"creation_date"`
	LastLogin    time.Time `json:"last_login"`
}
