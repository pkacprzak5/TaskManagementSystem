package store

import (
	"net/http"
)

type TaskService struct {
	store Store
}

func NewTaskService(store Store) *TaskService {
	return &TaskService{store: store}
}

func (s *TaskService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /tasks", WithJWTAuth(s.handleCreateTask, s.store))
	router.HandleFunc("GET /tasks/{id}", WithJWTAuth(s.handleCreateTask, s.store))
}
