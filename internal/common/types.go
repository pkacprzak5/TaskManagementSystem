package common

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	ID           int64     `json:"id" :"id"`
	Name         string    `json:"name" :"name"`
	Status       string    `json:"status" :"status"`
	ProjectID    int64     `json:"project_id" :"project_id"`
	AssignedToID int64     `json:"assigned_to_id" :"assigned_to_id"`
	CreatedAt    time.Time `json:"created_at" :"created_at"`
}

type User struct {
	ID        int64     `json:"id" :"id"`
	FirstName string    `json:"first_name" :"first_name"`
	LastName  string    `json:"last_name" :"last_name"`
	Email     string    `json:"email" :"email"`
	Password  string    `json:"password" :"password"`
	CreatedAt time.Time `json:"created_at" :"created_at"`
}
