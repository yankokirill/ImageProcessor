package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "hw/models"
)

var _ UserRepository = PostgresUserRepository{}

type UserRepository interface {
	AddUser(user *User) error
	ValidateUser(user *User) error
}

type PostgresUserRepository struct {
	pgPool *pgxpool.Pool
}

type InvalidCredentialsError struct {
	Message string
}

func (e *InvalidCredentialsError) Error() string {
	return e.Message
}

func NewInvalidCredentialsError(message string) error {
	return &InvalidCredentialsError{Message: message}
}

type UserExistsError struct{}

func (e *UserExistsError) Error() string {
	return "User already exists"
}

func NewUserExistsError() error {
	return &UserExistsError{}
}

func (r PostgresUserRepository) UserExists(login string) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)`
	err = r.pgPool.QueryRow(context.Background(), query, login).Scan(&exists)
	return
}

func (r PostgresUserRepository) AddUser(user *User) error {
	exists, err := r.UserExists(user.Login)
	if err != nil {
		return err
	}
	if exists {
		return NewUserExistsError()
	}

	query := `INSERT INTO users (user_id, login, password) VALUES ($1, $2, $3)`
	_, err = r.pgPool.Exec(context.Background(), query, user.ID, user.Login, user.Password)
	return err
}

func (r PostgresUserRepository) ValidateUser(user *User) error {
	var savedUser User
	query := `SELECT user_id, password FROM users WHERE login=$1`
	err := r.pgPool.QueryRow(context.Background(), query, user.Login).Scan(&savedUser.ID, &savedUser.Password)
	if err == pgx.ErrNoRows {
		return NewInvalidCredentialsError("Username doesn't exist")
	} else if err != nil {
		return err
	} else if savedUser.Password != user.Password {
		return NewInvalidCredentialsError("Wrong username or password")
	}
	user.ID = savedUser.ID
	return nil
}
