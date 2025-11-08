package store

import (
	"context"

	"backend/Infrastructure/store/dto"
	"backend/domain/model"
	"backend/domain/repository"

	"github.com/jinzhu/gorm"
)

// CategoryRepository implements category-specific persistence.
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a CategoryRepository.
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &CategoryRepository{db: db}
}

// ListCategories returns every category.
func (r *CategoryRepository) ListCategories(ctx context.Context) ([]model.Category, error) {
	var categoryDTOs []dto.Category
	if err := r.db.Find(&categoryDTOs).Error; err != nil {
		return nil, err
	}

	categories := make([]model.Category, 0, len(categoryDTOs))
	for _, c := range categoryDTOs {
		categories = append(categories, c.ToModel())
	}

	return categories, nil
}
