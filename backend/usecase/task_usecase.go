package usecase

import (
	"context"

	"backend/domain/model"
	"backend/domain/repository"
)

// TaskUseCase defines the business logic contract for tasks.
type TaskUseCase interface {
	ListTasks(ctx context.Context, categoryID *uint64) ([]model.Task, error)
	CreateTask(ctx context.Context, in model.Task) (*model.Task, error)
	UpdateTask(ctx context.Context, in model.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, id uint64) error
}

type taskUseCase struct {
	repo repository.TaskRepository
}

// NewTaskUseCase constructs a TaskUseCase implementation.
func NewTaskUseCase(repo repository.TaskRepository) TaskUseCase {
	return &taskUseCase{repo: repo}
}

// ListTasks returns all tasks.
func (uc *taskUseCase) ListTasks(ctx context.Context, categoryID *uint64) ([]model.Task, error) {
	return uc.repo.FindAll(ctx, categoryID)
}

// CreateTask creates and persists a new task.
func (uc *taskUseCase) CreateTask(ctx context.Context, in model.Task) (*model.Task, error) {
	return uc.repo.Create(ctx, in)
}

// UpdateTask updates an existing task.
func (uc *taskUseCase) UpdateTask(ctx context.Context, in model.UpdateTaskRequest) (*model.Task, error) {
	// 1. 既存データを取得
	task, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	// 2. nil でない項目のみ更新
	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Note != nil {
		task.Note = *in.Note
	}
	if in.Completed != nil {
		task.Completed = *in.Completed
	}
	if in.CategoryID != nil {
		task.CategoryID = *in.CategoryID
	}

	// 3. リポジトリ層に保存
	res, err := uc.repo.Update(ctx, *task)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteTask removes a task by id.
func (uc *taskUseCase) DeleteTask(ctx context.Context, id uint64) error {
	return uc.repo.Delete(ctx, id)
}
