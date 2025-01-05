package app

import (
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"net/http"
)

type TaskService struct {
	store common.Store
}

func NewTaskService(store common.Store) *TaskService {
	return &TaskService{store: store}
}

func (s *TaskService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /tasks", WithJWTAuth(s.handleCreateTask, s.store))
	router.HandleFunc("GET /tasks/{id}", WithJWTAuth(s.handleCreateTask, s.store))
}
