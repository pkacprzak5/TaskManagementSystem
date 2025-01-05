package common

import (
	"database/sql"
	"fmt"
	"time"
)

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)

	GetUserByID(id int) (*User, error)

	CreateTask(task *Task) (*Task, error)

	GetTask(id int) (*Task, error)

	UpdateTaskStatusByID(id int) (*Task, error)

	GetTasksAssignedToUser(id int) ([]*Task, error)
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
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) CreateTask(task *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, assignedToID) VALUES (?, ?, ?)",
		task.Name, task.Status, task.AssignedToID)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	task.ID = id
	task.CreatedAt = time.Now()
	return task, nil
}

func (s *Storage) GetTask(id int) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, assignedToID, createdAt FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Name, &t.Status, &t.AssignedToID, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func (s *Storage) UpdateTaskStatusByID(id int) (*Task, error) {
	task, err := s.GetTask(id)
	if err != nil {
		return &Task{}, fmt.Errorf("Task with id %d does not exist", id)
	}

	switch task.Status {
	case "TODO":
		task.Status = "IN_PROGRESS"
		break
	case "IN_PROGRESS":
		task.Status = "DONE"
		break
	case "DONE":
		return nil, fmt.Errorf("Task with id %d is already done", id)
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

	_, err = s.db.Exec(query, id)
	if err != nil {
		return &Task{}, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

func (s *Storage) GetTasksAssignedToUser(id int) ([]*Task, error) {
	query := "SELECT id, name, status, assignedToID, createdAt FROM tasks WHERE assignedToID = ?"

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks assigned to user with id %d: %w", id, err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Status, &task.AssignedToID, &task.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks assigned to user with id %d", id)
	}

	return tasks, nil
}
