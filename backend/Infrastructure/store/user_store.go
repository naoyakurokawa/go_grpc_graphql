package store

import (
	"context"

	"backend/Infrastructure/store/dto"
	"backend/domain/model"
	"backend/domain/repository"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user dto.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	result := user.ToModel()
	return &result, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user dto.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	result := user.ToModel()
	return &result, nil
}

func (r *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	d := dto.UserFromModel(user)
	if err := r.db.Create(&d).Error; err != nil {
		return nil, err
	}
	result := d.ToModel()
	return &result, nil
}
