package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	. "hw/messaging"
	. "hw/models"
	_ "hw/server/docs"
	. "hw/storage"
	"net/http"
	"strings"
)

type Server struct {
	storage Storage
	broker  Producer
}

type Response struct {
	Data  *Task
	Error string
	Code  int
}

func NewServer(storage Storage, broker Producer) *Server {
	return &Server{storage, broker}
}

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

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := s.storage.GetUserBySession(token)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func parseUserRequest(r *http.Request) (*User, error) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	if user.Login == "" || user.Password == "" {
		return nil, errors.New("invalid request")
	}
	return &user, nil
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
// @Failure 500 {string} string "Failed to store user"
// @Router /register [post]
func (s *Server) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	newUser, err := parseUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newUser.ID = uuid.New()
	if err := s.storage.AddUser(newUser); err != nil {
		if _, ok := err.(*UserExistsError); ok {
			http.Error(w, "User already exists", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to store user", http.StatusInternalServerError)
		}
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
// @Failure 400 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (s *Server) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := parseUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := s.storage.Login(user)
	if err != nil {
		if _, ok := err.(*InvalidCredentialsError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	sendJSON(w, "token", token)
}

func (s *Server) getTaskInfo(r *http.Request) Response {
	strTaskID := chi.URLParam(r, "task_id")
	taskID, err := uuid.Parse(strTaskID)
	if err != nil {
		return Response{nil, err.Error(), http.StatusBadRequest}
	}
	task, err := s.storage.GetTask(taskID)
	if err != nil {
		if _, ok := err.(*TaskNotFoundError); ok {
			return Response{nil, "Task not found", http.StatusNotFound}
		} else {
			return Response{nil, "Internal Server Error", http.StatusInternalServerError}
		}
	}
	userID := r.Context().Value("user_id").(uuid.UUID)
	if task.UserID != userID {
		return Response{nil, "Forbidden: You are not the owner of this task", http.StatusForbidden}
	}
	return Response{Data: &task}
}

// getStatusHandler retrieves the status of a task.
// @Summary GetTask task status
// @Description Retrieves the current status of the task by its ID.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string "Task Status"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /status/{task_id} [get]
func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := s.getTaskInfo(r)
	if response.Error != "" {
		http.Error(w, response.Error, response.Code)
		return
	}
	sendJSON(w, "status", response.Data.Status)
}

// getResultHandler retrieves the result of a task.
// @Summary GetTask task result
// @Description Retrieves the current result of the task by its ID.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string "Task Result"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /result/{task_id} [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	response := s.getTaskInfo(r)
	if response.Error != "" {
		http.Error(w, response.Error, response.Code)
		return
	}
	sendJSON(w, "result", response.Data.Result)
}

func (s *Server) createTask(r *http.Request) Response {
	task := &Task{
		ID:     uuid.New(),
		UserID: r.Context().Value("user_id").(uuid.UUID),
		Status: "in_progress",
	}
	if err := json.NewDecoder(r.Body).Decode(&task.Payload); err != nil {
		return Response{nil, "Invalid request", http.StatusBadRequest}
	}

	if err := s.storage.AddTask(task); err != nil {
		return Response{nil, "Failed to add task", http.StatusInternalServerError}
	}

	err := s.broker.Publish(task)
	if err != nil {
		return Response{nil, "Failed to enqueue task", http.StatusInternalServerError}
	}
	return Response{Data: task}
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
	response := s.createTask(r)
	if response.Error != "" {
		http.Error(w, response.Error, response.Code)
		return
	}
	w.WriteHeader(http.StatusCreated)
	sendJSON(w, "task_id", response.Data.ID.String())
}

func CreateAndRunServer(server *Server, addr string) error {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/", func(r chi.Router) {
		r.Post("/register", server.postRegisterHandler)
		r.Post("/login", server.postLoginHandler)
		r.Get("/status/{task_id}", server.AuthMiddleware(server.getStatusHandler))
		r.Get("/result/{task_id}", server.AuthMiddleware(server.getResultHandler))
		r.Post("/task", server.AuthMiddleware(server.postTaskHandler))
	})

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
