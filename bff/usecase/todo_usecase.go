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
	ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error)
	CreateSubTask(ctx context.Context, input model.NewSubTask) (*model.SubTask, error)
	ToggleSubTask(ctx context.Context, id uint64, completed bool) (*model.SubTask, error)
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

func (uc *todoUsecase) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error) {
	return uc.repo.ListTasks(ctx, filter)
}

func (uc *todoUsecase) CreateSubTask(ctx context.Context, input model.NewSubTask) (*model.SubTask, error) {
	return uc.repo.CreateSubTask(ctx, input)
}

func (uc *todoUsecase) ToggleSubTask(ctx context.Context, id uint64, completed bool) (*model.SubTask, error) {
	return uc.repo.ToggleSubTask(ctx, id, completed)
}
