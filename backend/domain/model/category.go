package model

import "time"

// Category represents a task category entity.
type Category struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
