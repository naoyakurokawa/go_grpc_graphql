package repository

import (
	"backend/domain/model"
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Create(ctx context.Context, user model.User) (*model.User, error)
}
