package usecase

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (uint64, error)
	GetUser(ctx context.Context, id uint64) (*model.User, error)
}

type authUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (uc *authUsecase) Login(ctx context.Context, email, password string) (uint64, error) {
	return uc.repo.Login(ctx, email, password)
}

func (uc *authUsecase) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return uc.repo.GetUser(ctx, id)
}
