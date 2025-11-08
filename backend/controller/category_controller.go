package controller

import (
	"context"

	"backend/usecase"

	pb "backend/pkg/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CategoryHandler bridges category gRPC requests with the use case layer.
type CategoryHandler struct {
	pb.UnimplementedCategoryServiceServer
	usecase usecase.CategoryUseCase
}

// NewCategoryHandler constructs a CategoryHandler.
func NewCategoryHandler(uc usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{usecase: uc}
}

// GetCategories returns all categories.
func (h *CategoryHandler) GetCategories(ctx context.Context, _ *emptypb.Empty) (*pb.CategoryList, error) {
	categories, err := h.usecase.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	pbCategories := make([]*pb.Category, 0, len(categories))
	for _, c := range categories {
		category := c
		pbCategories = append(pbCategories, &pb.Category{
			Id:   category.ID,
			Name: category.Name,
		})
	}

	return &pb.CategoryList{Categories: pbCategories}, nil
}
