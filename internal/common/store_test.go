package common

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdUser, err := store.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdUser.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	mockUser := &User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "securepassword",
		CreatedAt: time.Now(),
	}
	mockUser2 := &User{
		ID:        2,
		FirstName: "James",
		LastName:  "Doe",
		Email:     "james.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
			AddRow(mockUser.ID, mockUser.FirstName, mockUser.LastName, mockUser.Email, mockUser.Password, mockUser.CreatedAt))

	user, err := store.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	assert.NotEqual(t, mockUser2, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	task := &Task{
		Name:         "Sample Task",
		Status:       "TODO",
		AssignedToID: 1,
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.Name, task.Status, task.AssignedToID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdTask, err := store.CreateTask(task)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdTask.ID)
	assert.NotZero(t, createdTask.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTask(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	mockTask := &Task{
		ID:           1,
		Name:         "Sample Task",
		Status:       "TODO",
		AssignedToID: 1,
		CreatedAt:    time.Now(),
	}

	mock.ExpectQuery("SELECT id, name, status, assignedToID, createdAt FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "assignedToID", "createdAt"}).
			AddRow(mockTask.ID, mockTask.Name, mockTask.Status, mockTask.AssignedToID, mockTask.CreatedAt))

	task, err := store.GetTask(1)
	assert.NoError(t, err)
	assert.Equal(t, mockTask, task)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTaskStatusByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	mockTask := &Task{
		ID:           1,
		Name:         "Sample Task",
		Status:       "TODO",
		AssignedToID: 1,
		CreatedAt:    time.Now(),
	}

	mock.ExpectQuery("SELECT id, name, status, assignedToID, createdAt FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status", "assignedToID", "createdAt"}).
			AddRow(mockTask.ID, mockTask.Name, mockTask.Status, mockTask.AssignedToID, mockTask.CreatedAt))

	mock.ExpectExec("UPDATE tasks").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedTask, err := store.UpdateTaskStatusByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "IN_PROGRESS", updatedTask.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasksAssignedToUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	store := NewStore(db)

	mockTasks := []*Task{
		{ID: 1, Name: "Task 1", Status: "TODO", AssignedToID: 1, CreatedAt: time.Now()},
		{ID: 2, Name: "Task 2", Status: "IN_PROGRESS", AssignedToID: 1, CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "status", "assignedToID", "createdAt"})
	for _, task := range mockTasks {
		rows.AddRow(task.ID, task.Name, task.Status, task.AssignedToID, task.CreatedAt)
	}

	mock.ExpectQuery("SELECT id, name, status, assignedToID, createdAt FROM tasks WHERE assignedToID = ?").
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := store.GetTasksAssignedToUser(1)
	assert.NoError(t, err)
	assert.Equal(t, mockTasks, tasks)
	assert.NoError(t, mock.ExpectationsWereMet())
}
