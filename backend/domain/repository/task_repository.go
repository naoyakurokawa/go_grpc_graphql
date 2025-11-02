package repository

import (
	"backend/domain/model"
	"context"
)

// TaskRepository defines the contract for task persistence operations.
type TaskRepository interface {
	FindAll(ctx context.Context) ([]model.Task, error)
	FindByID(ctx context.Context, id uint64) (*model.Task, error)
	Create(ctx context.Context, in model.Task) (*model.Task, error)
	Update(ctx context.Context, in model.Task) (*model.Task, error)
	Delete(ctx context.Context, id uint64) error
}
