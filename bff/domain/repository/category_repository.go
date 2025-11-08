package repository

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
)

// CategoryRepository defines persistence operations for categories.
type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]*model.Category, error)
}
