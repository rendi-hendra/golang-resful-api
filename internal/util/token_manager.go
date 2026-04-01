package util

import "github.com/rendi-hendra/resful-api/internal/model"

// TokenManager defines the contract for token operations.
type TokenManager interface {
	CreateAccessToken(auth *model.Auth) (string, error)
	CreateRefreshToken(auth *model.Auth) (string, error)
	ParseToken(jwtToken string, expectedType string) (*model.Auth, error)
}
