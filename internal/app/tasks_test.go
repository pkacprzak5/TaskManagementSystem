package app

import (
	"bytes"
	"encoding/json"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(u *common.User) (*common.User, error) {
	args := m.Called(u)
	return args.Get(0).(*common.User), args.Error(1)
}

func (m *MockStore) GetUserByID(id int) (*common.User, error) {
	args := m.Called(id)
	return args.Get(0).(*common.User), args.Error(1)
}

func (m *MockStore) CreateTask(task *common.Task) (*common.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*common.Task), args.Error(1)
}

func (m *MockStore) GetTask(id int) (*common.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*common.Task), args.Error(1)
}

func (m *MockStore) UpdateTaskStatusByID(id int) (*common.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*common.Task), args.Error(1)
}

func (m *MockStore) GetTasksAssignedToUser(id int) ([]*common.Task, error) {
	args := m.Called(id)
	return args.Get(0).([]*common.Task), args.Error(1)
}

func TestHandleCreateTask(t *testing.T) {
	mockStore := new(MockStore)
	taskService := NewTaskService(mockStore)

	taskPayload := &common.Task{
		Name:         "Test Task",
		Status:       "TODO",
		AssignedToID: 1,
	}
	mockStore.On("CreateTask", taskPayload).Return(taskPayload, nil)

	requestBody, _ := json.Marshal(taskPayload)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	taskService.handleCreateTask(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdTask common.Task
	err := json.NewDecoder(resp.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, taskPayload.Name, createdTask.Name)
	mockStore.AssertExpectations(t)
}
