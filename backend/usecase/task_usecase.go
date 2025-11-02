package usecase

import (
	"context"
	"time"

	"backend/domain/model"
	"backend/domain/repository"
)

// TaskUseCase defines the business logic contract for tasks.
type TaskUseCase interface {
	ListTasks(ctx context.Context) ([]model.Task, error)
	CreateTask(ctx context.Context, in model.Task) (*model.Task, error)
	UpdateTask(ctx context.Context, id, title, note string, completed int32) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskUseCase struct {
	repo repository.TaskRepository
}

// NewTaskUseCase constructs a TaskUseCase implementation.
func NewTaskUseCase(repo repository.TaskRepository) TaskUseCase {
	return &taskUseCase{repo: repo}
}

// ListTasks returns all tasks.
func (uc *taskUseCase) ListTasks(ctx context.Context) ([]model.Task, error) {
	return uc.repo.FindAll(ctx)
}

// CreateTask creates and persists a new task.
func (uc *taskUseCase) CreateTask(ctx context.Context, in model.Task) (*model.Task, error) {
	return uc.repo.Create(ctx, in)
}

// UpdateTask updates an existing task.
func (uc *taskUseCase) UpdateTask(ctx context.Context, id, title, note string, completed int32) (*model.Task, error) {
	task, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Note = note
	task.Completed = completed
	task.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask removes a task by id.
func (uc *taskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
