package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "hw/models"
	"log"
	"os"
)

var _ TaskRepository = PostgresTaskRepository{}

type TaskRepository interface {
	GetTask(id uuid.UUID) (Task, error)
	AddTask(task *Task) error
	UpdateTaskStatus(id uuid.UUID, status, result string)
}

type PostgresTaskRepository struct {
	pgPool *pgxpool.Pool
}

type TaskNotFoundError struct{}

func (e *TaskNotFoundError) Error() string {
	return "Task not found"
}

func NewTaskNotFoundError() error {
	return &TaskNotFoundError{}
}

func NewPostgresTaskRepo(connString string) PostgresTaskRepository {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	return PostgresTaskRepository{pool}
}

func (r PostgresTaskRepository) GetTask(id uuid.UUID) (Task, error) {
	var task Task
	query := `SELECT task_id, user_id, status, result FROM tasks WHERE task_id=$1`
	err := r.pgPool.QueryRow(context.Background(), query, id).Scan(&task.ID, &task.UserID, &task.Status, &task.Result)
	if err == pgx.ErrNoRows {
		return Task{}, NewTaskNotFoundError()
	}
	return task, err
}

func (r PostgresTaskRepository) AddTask(task *Task) error {
	payloadData, err := json.Marshal(task.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	query := `INSERT INTO tasks (task_id, user_id, payload, status, result) VALUES ($1, $2, $3, $4, $5)`
	_, err = r.pgPool.Exec(context.Background(), query, task.ID, task.UserID, payloadData, task.Status, task.Result)
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}
	return nil
}

func (r PostgresTaskRepository) UpdateTaskStatus(id uuid.UUID, status, result string) {
	query := `UPDATE tasks SET status=$1, result=$2 WHERE task_id=$3`
	_, err := r.pgPool.Exec(context.Background(), query, status, result, id)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error updating task:", err)
	}
}
