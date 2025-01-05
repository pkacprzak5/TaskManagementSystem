package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"github.com/pkacprzak5/TaskManagementSystem/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store common.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("failed to authenticate token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("failed to authenticate token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(int)

		_, err = store.GetUserByID(userID)
		if err != nil {
			fmt.Println("failed to get user")
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}
	return ""
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, common.ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func validateJWT(token string) (*jwt.Token, error) {
	secret := common.Envs.JWTSecret

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func HashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
