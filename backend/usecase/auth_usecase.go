package usecase

import (
	"context"

	"backend/domain/model"
	"backend/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (*model.User, error)
	GetUser(ctx context.Context, id uint64) (*model.User, error)
}

type authUseCase struct {
	repo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUseCase) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return uc.repo.FindByID(ctx, id)
}
