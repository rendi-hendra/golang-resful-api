package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rendi-hendra/resful-api/internal/model"
	"github.com/rendi-hendra/resful-api/internal/util"
	"github.com/sirupsen/logrus"
)

func NewAuth(log *logrus.Logger, tokenUtil util.TokenManager) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		token, err := ExtractToken(header)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		log.Debugf("Authorization : %s", token)

		auth, err := tokenUtil.ParseToken(token, "access")
		if err != nil {
			log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func ExtractToken(header string) (string, error) {
	if header == "" {
		return "", fiber.ErrUnauthorized
	}

	const prefix = "Bearer "
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return "", fiber.ErrUnauthorized
	}

	return header[len(prefix):], nil
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
