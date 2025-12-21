package controller

import (
	"context"
	"errors"

	"backend/usecase"

	pb "backend/pkg/pb"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const invalidCredentialsMessage = "メールアドレスまたはパスワードが正しくありません"

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
		if isInvalidCredentialError(err) {
			return nil, status.Error(codes.Unauthenticated, invalidCredentialsMessage)
		}
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
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

func isInvalidCredentialError(err error) bool {
	return gorm.IsRecordNotFoundError(err) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
}
