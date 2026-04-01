package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rendi-hendra/resful-api/internal/model"
)

type TokenUtil struct {
	SecretKey string
}

func NewTokenUtil(secretKey string) TokenManager {
	return &TokenUtil{
		SecretKey: secretKey,
	}
}

func (t TokenUtil) CreateAccessToken(auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   auth.ID,
		"type": "access",
		"exp":  time.Now().Add(time.Minute * 1).Unix(),
	})

	jwtAccessToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtAccessToken, nil
}

func (t TokenUtil) CreateRefreshToken(auth *model.Auth) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   auth.ID,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 7).Unix(),
	})

	jwtRefreshToken, err := refreshToken.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtRefreshToken, nil

}

func (t TokenUtil) ParseToken(jwtToken string, expectedType string) (*model.Auth, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(t.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	// validate token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fiber.ErrUnauthorized
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	auth := &model.Auth{ID: id}
	return auth, nil
}
