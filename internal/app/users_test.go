package app

import (
	"errors"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"testing"
	"time"
)

type mockStore struct {
	CreateUserFunc             func(u *common.User) (*common.User, error)
	GetUserByIDFunc            func(id int) (*common.User, error)
	CreateTaskFunc             func(task *common.Task) (*common.Task, error)
	GetTaskFunc                func(id int) (*common.Task, error)
	UpdateTaskStatusByIDFunc   func(id int) (*common.Task, error)
	GetTasksAssignedToUserFunc func(id int) ([]*common.Task, error)
}

func (m *mockStore) CreateUser(u *common.User) (*common.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(u)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStore) GetUserByID(id int) (*common.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStore) CreateTask(task *common.Task) (*common.Task, error) {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(task)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStore) GetTask(id int) (*common.Task, error) {
	if m.GetTaskFunc != nil {
		return m.GetTaskFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStore) UpdateTaskStatusByID(id int) (*common.Task, error) {
	if m.UpdateTaskStatusByIDFunc != nil {
		return m.UpdateTaskStatusByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStore) GetTasksAssignedToUser(id int) ([]*common.Task, error) {
	if m.GetTasksAssignedToUserFunc != nil {
		return m.GetTasksAssignedToUserFunc(id)
	}
	return nil, errors.New("not implemented")
}

func TestCreateUser_Success(t *testing.T) {
	mock := &mockStore{
		CreateUserFunc: func(u *common.User) (*common.User, error) {
			u.ID = 1
			u.CreatedAt = time.Now()
			return u, nil
		},
	}

	user := &common.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "securepassword",
	}

	result, err := mock.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID == 0 {
		t.Fatalf("expected user ID to be set, got %d", result.ID)
	}
}

func TestCreateTask_Success(t *testing.T) {
	mock := &mockStore{
		CreateTaskFunc: func(task *common.Task) (*common.Task, error) {
			task.ID = 1
			task.CreatedAt = time.Now()
			return task, nil
		},
	}

	task := &common.Task{
		Name:         "Sample Task",
		Status:       "TODO",
		AssignedToID: 1,
	}

	result, err := mock.CreateTask(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID == 0 {
		t.Fatalf("expected task ID to be set, got %d", result.ID)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	mock := &mockStore{
		GetUserByIDFunc: func(id int) (*common.User, error) {
			return &common.User{
				ID:        int64(id),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	user, err := mock.GetUserByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected user ID to be 1, got %d", user.ID)
	}
}
