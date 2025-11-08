package controller

import (
	"context"
	"log"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/usecase"
)

// CategoryController orchestrates category related operations.
type CategoryController struct {
	usecase usecase.CategoryUsecase
}

// NewCategoryController constructs a CategoryController instance.
func NewCategoryController(uc usecase.CategoryUsecase) *CategoryController {
	return &CategoryController{usecase: uc}
}

func (c *CategoryController) ListCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := c.usecase.ListCategories(ctx)
	if err != nil {
		log.Printf("failed to fetch categories: %v", err)
		return nil, err
	}

	return categories, nil
}
