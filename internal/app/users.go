package app

import (
	"encoding/json"
	"errors"
	"github.com/pkacprzak5/TaskManagementSystem/internal/auth"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"github.com/pkacprzak5/TaskManagementSystem/pkg/utils"
	"io"
	"net/http"
)

type UsersService struct {
	store common.Store
}

func NewUsersService(store common.Store) *UsersService {
	return &UsersService{store}
}

func (s *UsersService) RegusterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", s.handleUserRegister)
}

func (s *UsersService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload common.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if err := validateUserPayload(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	payload.Password = hashedPassword

	user, err := s.store.CreateUser(&payload)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)

}

func validateUserPayload(u common.User) error {
	if u.Email == "" {
		return errors.New("Email is required")
	}

	if u.Password == "" {
		return errors.New("Password is required")
	}

	if u.FirstName == "" {
		return errors.New("FirstName is required")
	}

	if u.LastName == "" {
		return errors.New("LastName is required")
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(common.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
