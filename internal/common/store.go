package common

import (
	"database/sql"
	"fmt"
)

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)

	GetUserByID(id int) (*User, error)

	CreateTask(task *Task) (*Task, error)

	GetTask(id int) (*Task, error)

	UpdateTaskStatusByID(id int) error
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)",
		u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (s *Storage) GetUserByID(id int) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) CreateTask(task *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, project_id, assigned_to) VALUES (?, ?, ?, ?)",
		task.Name, task.Status, task.ProjectID, task.AssignedToID)
	if err != nil {
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	task.ID = id
	return task, nil
}

func (s *Storage) GetTask(id int) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, project_id, assigned_to FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID)
	return &t, err
}

func (s *Storage) UpdateTaskStatusByID(id int) error {
	_, err := s.GetTask(id)
	if err != nil {
		return fmt.Errorf("Task with id %d does not exist", id)
	}

	query := `
		UPDATE tasks
		SET status = CASE status
			WHEN 'TODO' THEN 'IN_PROGRESS'
			WHEN 'IN_PROGRESS' THEN 'IN_TESTING'
			WHEN 'IN_TESTING' THEN 'DONE'
			ELSE status
		END
		WHERE id = ? AND status != 'DONE'; -- Prevent updating if already DONE
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task status not updated, possibly already DONE or task not found")
	}

	return nil
}
