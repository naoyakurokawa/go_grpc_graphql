package store

import (
	"context"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
	pb "github.com/naoyakurokawa/go_grpc_graphql/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ repository.AuthRepository = (*AuthStore)(nil)

type AuthStore struct {
	client pb.AuthServiceClient
}

func NewAuthStore(client pb.AuthServiceClient) repository.AuthRepository {
	return &AuthStore{client: client}
}

func (s *AuthStore) Login(ctx context.Context, email, password string) (uint64, error) {
	res, err := s.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return res.GetUserId(), nil
}

func (s *AuthStore) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        user.GetId(),
		Email:     user.GetEmail(),
		CreatedAt: formatUserTimestamp(user.GetCreatedAt()),
		UpdatedAt: formatUserTimestamp(user.GetUpdatedAt()),
	}, nil
}

func formatUserTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format("2006-01-02 15:04:05")
}
