package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rendi-hendra/resful-api/internal/model"
	"github.com/rendi-hendra/resful-api/internal/usecase"
	"github.com/rendi-hendra/resful-api/internal/util"
)

func NewAuth(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := tokenUtil.ParseToken(request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
