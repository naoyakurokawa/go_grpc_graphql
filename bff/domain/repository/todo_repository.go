package repository

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
)

// TodoRepository defines the persistence contract for todo related operations.
type TodoRepository interface {
	CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error)
	UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error)
	DeleteTask(ctx context.Context, id uint64) (bool, error)
	ListTasks(ctx context.Context, categoryID *uint64) ([]*model.Task, error)
}
