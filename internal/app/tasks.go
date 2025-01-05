package app

import (
	"encoding/json"
	"errors"
	"github.com/pkacprzak5/TaskManagementSystem/internal/auth"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"github.com/pkacprzak5/TaskManagementSystem/pkg/utils"
	"io"
	"net/http"
	"strconv"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TaskService struct {
	store common.Store
}

func NewTaskService(store common.Store) *TaskService {
	return &TaskService{store: store}
}

func (s *TaskService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /tasks", auth.WithJWTAuth(s.handleCreateTask, s.store))
	router.HandleFunc("GET /tasks/{id}", auth.WithJWTAuth(s.handleGetTask, s.store))
}

func (s *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload common.Task
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if err := validateTaskPayload(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.store.CreateTask(&payload)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, task)

}

func (s *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	task, err := s.store.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

func validateTaskPayload(task *common.Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
