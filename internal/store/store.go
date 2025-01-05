package store

import "database/sql"

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)

	GetUserByID(id string) (*User, error)

	CreateTask(task *Task) (*Task, error)

	GetTask(id string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{db: db}
}
