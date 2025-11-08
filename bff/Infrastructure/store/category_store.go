package store

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
	pb "github.com/naoyakurokawa/go_grpc_graphql/pkg/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ repository.CategoryRepository = (*CategoryStore)(nil)

// CategoryStore implements CategoryRepository via gRPC.
type CategoryStore struct {
	client pb.CategoryServiceClient
}

// NewCategoryStore creates a CategoryStore.
func NewCategoryStore(client pb.CategoryServiceClient) repository.CategoryRepository {
	return &CategoryStore{client: client}
}

func (s *CategoryStore) ListCategories(ctx context.Context) ([]*model.Category, error) {
	res, err := s.client.GetCategories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	categories := make([]*model.Category, 0, len(res.Categories))
	for _, c := range res.Categories {
		categories = append(categories, &model.Category{
			ID:   c.GetId(),
			Name: c.GetName(),
		})
	}

	return categories, nil
}
