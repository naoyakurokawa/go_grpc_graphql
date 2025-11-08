package repository

import (
	"backend/domain/model"
	"context"
)

type SubTaskRepository interface {
	ListByTaskID(ctx context.Context, taskID uint64) ([]model.SubTask, error)
	Create(ctx context.Context, in model.SubTask) (*model.SubTask, error)
	Update(ctx context.Context, in model.SubTask) (*model.SubTask, error)
	FindByID(ctx context.Context, id uint64) (*model.SubTask, error)
}
