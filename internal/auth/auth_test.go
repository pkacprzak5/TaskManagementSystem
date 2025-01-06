package auth_test

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/pkacprzak5/TaskManagementSystem/internal/auth"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStore struct {
	users map[int]*common.User
}

func (m *MockStore) CreateUser(u *common.User) (*common.User, error) { return nil, nil }
func (m *MockStore) GetUserByID(id int) (*common.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
func (m *MockStore) CreateTask(task *common.Task) (*common.Task, error)    { return nil, nil }
func (m *MockStore) GetTask(id int) (*common.Task, error)                  { return nil, nil }
func (m *MockStore) UpdateTaskStatusByID(id int) (*common.Task, error)     { return nil, nil }
func (m *MockStore) GetTasksAssignedToUser(id int) ([]*common.Task, error) { return nil, nil }

func TestWithJWTAuth_Success(t *testing.T) {
	store := &MockStore{
		users: map[int]*common.User{
			1: {ID: 1, FirstName: "John", LastName: "Doe"},
		},
	}

	handler := auth.WithJWTAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, store)

	secret := []byte("testsecret")
	token, _ := auth.CreateJWT(secret, 1)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token)

	common.Envs.JWTSecret = string(secret)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestWithJWTAuth_PermissionDenied(t *testing.T) {
	store := &MockStore{
		users: map[int]*common.User{},
	}

	handler := auth.WithJWTAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, store)

	secret := []byte("testsecret")
	token, _ := auth.CreateJWT(secret, 1)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token)

	common.Envs.JWTSecret = string(secret)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetUserIDFromRequest_Success(t *testing.T) {
	secret := []byte("testsecret")
	token, _ := auth.CreateJWT(secret, 1)

	common.Envs.JWTSecret = string(secret)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token)

	userID, err := auth.GetUserIDFromRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

func TestGetUserIDFromRequest_InvalidToken(t *testing.T) {
	common.Envs.JWTSecret = "wrongsecret"
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "invalidtoken")

	userID, err := auth.GetUserIDFromRequest(req)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}

func TestHashedPassword(t *testing.T) {
	password := "password123"
	hashed, err := auth.HashedPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	assert.NoError(t, err)
}

func TestCreateJWT(t *testing.T) {
	secret := []byte("testsecret")
	userID := int64(123)
	token, err := auth.CreateJWT(secret, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(userID), claims["userID"])
	assert.NotNil(t, claims["expiresAt"])
}

func TestGetTokenFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/?token=querytoken", nil)
	req.Header.Set("Authorization", "headertoken")

	token := auth.GetTokenFromRequest(req)
	assert.Equal(t, "headertoken", token)

	req = httptest.NewRequest("GET", "/?token=querytoken", nil)
	token = auth.GetTokenFromRequest(req)
	assert.Equal(t, "querytoken", token)

	req = httptest.NewRequest("GET", "/", nil)
	token = auth.GetTokenFromRequest(req)
	assert.Equal(t, "", token)
}
