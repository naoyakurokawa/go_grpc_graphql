package model

import "time"

// Task represents a todo task entity.
type Task struct {
	ID        uint64
	Title     string
	Note      string
	Completed int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateTaskRequest struct {
	ID        uint64
	Title     *string
	Note      *string
	Completed *int32
}
