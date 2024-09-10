package ram_storage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	. "github.com/yankokirill/ImageProcessor/models"
	"os"
)

type Storage interface {
	Get(key uuid.UUID) (Task, error)

	AddTask(task Task) error

	AddUser(user User) error

	Login(user User) (uuid.UUID, error)

	SessionExists(token uuid.UUID) bool

	UpdateTaskStatus(id uuid.UUID, status, result string)
}

type RamStorage struct {
	taskData    map[uuid.UUID]Task
	userData    map[string]User
	sessionData map[uuid.UUID]Session
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		taskData:    make(map[uuid.UUID]Task),
		userData:    make(map[string]User),
		sessionData: make(map[uuid.UUID]Session),
	}
}

func (rs *RamStorage) Get(key uuid.UUID) (Task, error) {
	value, exists := rs.taskData[key]
	if !exists {
		return Task{}, errors.New("key not found")
	}
	return value, nil
}

func (rs *RamStorage) AddTask(task Task) error {
	if _, exists := rs.taskData[task.ID]; exists {
		return errors.New("key already exists")
	}
	rs.taskData[task.ID] = task
	return nil
}

func (rs *RamStorage) AddUser(user User) error {
	if _, exists := rs.userData[user.Login]; exists {
		return errors.New("key already exists")
	}
	rs.userData[user.Login] = user
	return nil
}

func (rs *RamStorage) AddSession(session Session) error {
	if _, exists := rs.sessionData[session.SessionID]; exists {
		return errors.New("key already exists")
	}
	rs.sessionData[session.SessionID] = session
	return nil
}

func (rs *RamStorage) Login(user User) (uuid.UUID, error) {
	savedUser, ok := rs.userData[user.Login]
	if !ok {
		return uuid.UUID{}, errors.New("username doesn't exist")
	}
	if savedUser.Password != user.Password {
		return uuid.UUID{}, errors.New("wrong username or password")
	}
	token := uuid.New()
	session := Session{UserID: user.ID, SessionID: token}
	rs.sessionData[token] = session
	return token, nil
}

func (rs *RamStorage) SessionExists(token uuid.UUID) bool {
	_, ok := rs.sessionData[token]
	return ok
}

func (rs *RamStorage) UpdateTaskStatus(id uuid.UUID, status, result string) {
	value, ok := rs.taskData[id]
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "updating nonexistent task")
		os.Exit(1)
	}

	value.Status = status
	value.Result = result
	rs.taskData[id] = value
}
