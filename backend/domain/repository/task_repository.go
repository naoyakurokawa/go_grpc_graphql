package repository

import (
	"backend/domain/model"
	"context"
	"time"
)

// TaskRepository defines the contract for task persistence operations.
type TaskRepository interface {
	FindAll(ctx context.Context, filter TaskFilter) ([]model.Task, error)
	FindByID(ctx context.Context, id uint64) (*model.Task, error)
	Create(ctx context.Context, in model.Task) (*model.Task, error)
	Update(ctx context.Context, in model.Task) (*model.Task, error)
	Delete(ctx context.Context, id uint64) error
}

type TaskFilter struct {
	CategoryID     *uint64
	DueDateFrom    *time.Time
	DueDateTo      *time.Time
	IncompleteOnly *bool
	UserID         *uint64
}
