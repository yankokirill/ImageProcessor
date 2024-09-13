package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	. "hw/models"
	"log"
)

var _ Storage = &DatabaseStorage{}

type Storage interface {
	GetTask(id uuid.UUID) (Task, error)
	AddTask(task *Task) error
	UpdateTaskStatus(id uuid.UUID, status, result string)

	AddUser(user *User) error
	Login(user *User) (string, error)

	GetUserBySession(token string) (uuid.UUID, error)
}

type DatabaseStorage struct {
	PostgresTaskRepository
	PostgresUserRepository
	RedisSessionRepository
}

func NewDatabaseStorage(connString, redisAddr, jwtSecret string) *DatabaseStorage {
	taskRepo := NewPostgresTaskRepo(connString)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	return &DatabaseStorage{
		PostgresTaskRepository{taskRepo.pgPool},
		PostgresUserRepository{taskRepo.pgPool},
		RedisSessionRepository{
			redisClient: rdb,
			jwtSecret:   []byte(jwtSecret),
		},
	}
}

func (ds *DatabaseStorage) Login(user *User) (string, error) {
	err := ds.ValidateUser(user)
	if err != nil {
		return "", err
	}
	return ds.AddSession(user.ID.String())
}
