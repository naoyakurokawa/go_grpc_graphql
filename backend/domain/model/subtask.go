package model

import "time"

// SubTask represents a sub task associated with a parent task.
type SubTask struct {
	ID          uint64
	TaskID      uint64
	Title       string
	Note        string
	Completed   int32
	CompletedAt *time.Time
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
