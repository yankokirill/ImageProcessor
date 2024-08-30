package models

import "github.com/google/uuid"

type Storage interface {
	Get(key uuid.UUID) (Task, error)

	AddTask(task Task) error

	AddUser(user User) error

	Login(user User) (uuid.UUID, error)

	SessionExists(token uuid.UUID) bool

	UpdateTaskStatus(id uuid.UUID, status, result string)
}
