package usecase

import (
	"context"

	"backend/domain/model"
	"backend/domain/repository"
)

// CategoryUseCase defines category-specific business logic.
type CategoryUseCase interface {
	ListCategories(ctx context.Context) ([]model.Category, error)
}

type categoryUseCase struct {
	repo repository.CategoryRepository
}

// NewCategoryUseCase constructs a CategoryUseCase.
func NewCategoryUseCase(repo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{repo: repo}
}

// ListCategories returns every category.
func (uc *categoryUseCase) ListCategories(ctx context.Context) ([]model.Category, error) {
	return uc.repo.ListCategories(ctx)
}
