package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	. "models"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "server/docs"
)

// Server defines the HTTP server for handling task-related requests.
type Server struct {
	storage Storage
	ch      *amqp.Channel
}

type Response struct {
	Data  *Task
	Error string
	Code  int
}

func NewServer(storage Storage, ch *amqp.Channel) *Server {
	return &Server{storage, ch}
}

func (s *Server) sendTask(task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = s.ch.Publish(
		"",
		"task_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}

// AuthMiddleware checks for a valid authorization token in the request header
func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := uuid.Parse(tokenStr)

		if err != nil || !s.storage.SessionExists(token) {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func parseUserRequest(r *http.Request) (user User, err error) {
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		return
	}
	if user.Login == "" || user.Password == "" {
		err = errors.New("invalid request")
	}
	return
}

func parseTaskRequest(r *http.Request) (uuid.UUID, error) {
	taskID := chi.URLParam(r, "task_id")
	return uuid.Parse(taskID)
}

func sendJSON(w http.ResponseWriter, key, value string) {
	response := map[string]string{key: value}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// postRegisterHandler registers a new user.
// @Summary Register a new user
// @Description Creates a new user account.
// @Tags user
// @Accept  json
// @Produce  json
// @Success 201 "User registered successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to store value"
// @Router /register [post]
func (s *Server) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	newUser, err := parseUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newUser.ID = uuid.New()
	if err := s.storage.AddUser(newUser); err != nil {
		http.Error(w, "Failed to store value", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// postLoginHandler logs in an existing user.
// @Summary Log in a user
// @Description Authenticates a user and returns a token.
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string "Token"
// @Failure 400 {string} string "Invalid request"
// @Router /login [post]
func (s *Server) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := parseUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := s.storage.Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sendJSON(w, "token", token.String())
}

func (s *Server) processRequest(taskID uuid.UUID) Response {
	value, err := s.storage.Get(taskID)
	if err != nil {
		return Response{nil, "Task not found", http.StatusNotFound}
	}
	return Response{Data: &value}
}

// getStatusHandler retrieves the status of a task.
// @Summary Get task status
// @Description Retrieves the current status of the task by its id.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string "Task Status"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Router /status/{task_id} [get]
func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := s.processRequest(taskID)
	if response.Error != "" {
		http.Error(w, response.Error, response.Code)
		return
	}
	sendJSON(w, "status", response.Data.Status)
}

// getResultHandler retrieves the result of a task.
// @Summary Get task result
// @Description Retrieves the current result of the task by its id.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string "Task Result"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Router /result/{task_id} [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := s.processRequest(taskID)
	if response.Error != "" {
		http.Error(w, response.Error, response.Code)
		return
	}
	sendJSON(w, "result", response.Data.Result)
}

// postTaskHandler handles task creation requests.
// @Summary Create a new task
// @Description Creates a new task, sends it to ImageProcessor and returns the task ID.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 201 {object} map[string]string "Task ID"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to add task"
// @Router /task [post]
func (s *Server) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var data ImageProcessorPayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	taskID := uuid.New()
	task := Task{
		ID:      taskID,
		Payload: data,
		Status:  "in_progress",
	}
	if err := s.storage.AddTask(task); err != nil {
		http.Error(w, "Failed to add task", http.StatusInternalServerError)
		return
	}

	err := s.sendTask(task)
	if err != nil {
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendJSON(w, "task_id", taskID.String())
}

// postCommitHandler updates the task status and result.
func (s *Server) postCommitHandler(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	data := struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.storage.UpdateTaskStatus(taskID, data.Status, data.Result)
}

// CreateAndRunServer initializes and starts the HTTP server
func CreateAndRunServer(server *Server, addr string) error {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/", func(r chi.Router) {
		r.Post("/register", server.postRegisterHandler)
		r.Post("/login", server.postLoginHandler)
		r.Get("/status/{task_id}", server.AuthMiddleware(server.getStatusHandler))
		r.Get("/result/{task_id}", server.AuthMiddleware(server.getResultHandler))
		r.Post("/task", server.AuthMiddleware(server.postTaskHandler))
		r.Post("/commit/{task_id}", server.postCommitHandler)
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
