package repository

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
)

type AuthRepository interface {
	Login(ctx context.Context, email, password string) (uint64, error)
	GetUser(ctx context.Context, id uint64) (*model.User, error)
}
