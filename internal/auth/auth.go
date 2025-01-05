package auth

import (
	"github.com/pkacprzak5/TaskManagementSystem/internal/store"
	"net/http"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
}
