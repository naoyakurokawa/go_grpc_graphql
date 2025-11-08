package repository

import (
	"backend/domain/model"
	"context"
)

// CategoryRepository defines persistence operations for categories.
type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]model.Category, error)
}
