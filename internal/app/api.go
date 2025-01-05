package app

import (
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"log"
	"net/http"
)

type APIServer struct {
	address string
	store   common.Store
}

func NewAPIServer(address string, store common.Store) *APIServer {
	return &APIServer{address: address, store: store}
}

func (s *APIServer) Serve() {
	router := http.NewServeMux()

	usersService := NewUsersService(s.store)
	usersService.RegusterRoutes(router)

	tasksService := NewTaskService(s.store)
	tasksService.RegisterRoutes(router)

	log.Println("Starting API server at", s.address)

	log.Fatal(http.ListenAndServe(s.address, router))
}
