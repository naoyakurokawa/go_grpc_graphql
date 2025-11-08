package usecase

import (
	"context"
	"time"

	"backend/domain/model"
	"backend/domain/repository"
)

type SubTaskUseCase interface {
	ListByTaskID(ctx context.Context, taskID uint64) ([]model.SubTask, error)
	Create(ctx context.Context, in model.SubTask) (*model.SubTask, error)
	ToggleCompletion(ctx context.Context, id uint64, completed bool) (*model.SubTask, error)
}

type subTaskUseCase struct {
	repo repository.SubTaskRepository
}

func NewSubTaskUseCase(repo repository.SubTaskRepository) SubTaskUseCase {
	return &subTaskUseCase{repo: repo}
}

func (uc *subTaskUseCase) ListByTaskID(ctx context.Context, taskID uint64) ([]model.SubTask, error) {
	return uc.repo.ListByTaskID(ctx, taskID)
}

func (uc *subTaskUseCase) Create(ctx context.Context, in model.SubTask) (*model.SubTask, error) {
	return uc.repo.Create(ctx, in)
}

func (uc *subTaskUseCase) ToggleCompletion(ctx context.Context, id uint64, completed bool) (*model.SubTask, error) {
	subTask, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if completed {
		now := time.Now()
		subTask.Completed = 1
		subTask.CompletedAt = &now
	} else {
		subTask.Completed = 0
		subTask.CompletedAt = nil
	}

	return uc.repo.Update(ctx, *subTask)
}
