package usecase

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
)

// TodoUsecase defines business logic for todo operations.
type TodoUsecase interface {
	CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error)
	UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error)
	DeleteTask(ctx context.Context, id uint64) (bool, error)
	ListTasks(ctx context.Context, categoryID *uint64) ([]*model.Task, error)
}

type todoUsecase struct {
	repo repository.TodoRepository
}

// NewTodoUsecase creates a TodoUsecase backed by the provided repository.
func NewTodoUsecase(repo repository.TodoRepository) TodoUsecase {
	return &todoUsecase{repo: repo}
}

func (uc *todoUsecase) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	return uc.repo.CreateTask(ctx, input)
}

func (uc *todoUsecase) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	return uc.repo.UpdateTask(ctx, input)
}

func (uc *todoUsecase) DeleteTask(ctx context.Context, id uint64) (bool, error) {
	return uc.repo.DeleteTask(ctx, id)
}

func (uc *todoUsecase) ListTasks(ctx context.Context, categoryID *uint64) ([]*model.Task, error) {
	return uc.repo.ListTasks(ctx, categoryID)
}
