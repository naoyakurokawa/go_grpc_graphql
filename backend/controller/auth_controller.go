package controller

import (
	"context"

	"backend/usecase"

	pb "backend/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	usecase usecase.AuthUseCase
}

func NewAuthController(uc usecase.AuthUseCase) *AuthController {
	return &AuthController{usecase: uc}
}

func (c *AuthController) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := c.usecase.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{UserId: user.ID}, nil
}

func (c *AuthController) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	user, err := c.usecase.GetUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
