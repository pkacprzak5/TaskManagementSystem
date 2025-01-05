package app

import (
	"context"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	address string
	store   common.Store
}

func NewAPIServer(address string, store common.Store) *APIServer {
	return &APIServer{address: address, store: store}
}

func (s *APIServer) Serve(ctx context.Context) error {
	ch := make(chan error, 1)
	router := http.NewServeMux()

	usersService := NewUsersService(s.store)
	usersService.RegusterRoutes(router)

	tasksService := NewTaskService(s.store)
	tasksService.RegisterRoutes(router)

	server := &http.Server{
		Addr:    s.address,
		Handler: router,
	}

	go func() {
		log.Println("Starting API server at", s.address)

		err := server.ListenAndServe()
		if err != nil {
			ch <- err
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		log.Println("Received shutdown signal, shutting down the server...")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
