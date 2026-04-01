package usecase

import (
	"context"

	"github.com/rendi-hendra/resful-api/internal/model"
)

// UserUseCase defines the contract for all User business logic operations.
type UserUseCase interface {
	Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.TokenResponse, error)
	Refresh(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error)
	Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
}
