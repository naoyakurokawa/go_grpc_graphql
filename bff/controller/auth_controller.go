package controller

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/usecase"
)

type AuthController struct {
	usecase usecase.AuthUsecase
}

func NewAuthController(uc usecase.AuthUsecase) *AuthController {
	return &AuthController{usecase: uc}
}

func (c *AuthController) Login(ctx context.Context, email, password string) (uint64, error) {
	return c.usecase.Login(ctx, email, password)
}

func (c *AuthController) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return c.usecase.GetUser(ctx, id)
}
