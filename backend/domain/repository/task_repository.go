package repository

import (
	"backend/domain/model"
	"context"
)

// TaskRepository defines the contract for task persistence operations.
type TaskRepository interface {
	FindAll(ctx context.Context) ([]model.Task, error)
	FindByID(ctx context.Context, id string) (*model.Task, error)
	Create(ctx context.Context, task model.Task) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
}
