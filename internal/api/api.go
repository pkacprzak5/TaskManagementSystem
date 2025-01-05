package api

import (
	"github.com/pkacprzak5/TaskManagementSystem/internal/store"
	"log"
	"net/http"
)

type APIServer struct {
	address string
	store   store.Store
}

func NewAPIServer(address string, store store.Store) *APIServer {
	return &APIServer{address: address, store: store}
}

func (s *APIServer) Serve() {
	router := http.NewServeMux()

	log.Println("Starting API server at", s.address)

	log.Fatal(http.ListenAndServe(s.address, router))
}
