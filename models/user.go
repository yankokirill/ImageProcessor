package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"user_id"`
	Login    string    `json:"username"`
	Password string    `json:"password"`
}
