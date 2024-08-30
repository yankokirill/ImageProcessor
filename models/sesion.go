package models

import "github.com/google/uuid"

type Session struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}
