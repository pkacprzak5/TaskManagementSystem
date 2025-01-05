package store

import (
	"net/http"
)

type UsersService struct {
	store Store
}

func NewUsersService(store Store) *UsersService {
	return &UsersService{store}
}

func (s *UsersService) RegusterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", s.handleUserRegister)
}
