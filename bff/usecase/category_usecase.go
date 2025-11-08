package usecase

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
)

// CategoryUsecase exposes category business logic.
type CategoryUsecase interface {
	ListCategories(ctx context.Context) ([]*model.Category, error)
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

// NewCategoryUsecase creates a CategoryUsecase backed by the provided repository.
func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (uc *categoryUsecase) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return uc.repo.ListCategories(ctx)
}
