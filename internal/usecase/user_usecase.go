package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rendi-hendra/resful-api/internal/entity"
	"github.com/rendi-hendra/resful-api/internal/model"
	"github.com/rendi-hendra/resful-api/internal/model/converter"
	"github.com/rendi-hendra/resful-api/internal/repository"
	"github.com/rendi-hendra/resful-api/internal/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
	Mailer         *util.Mailer
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, tokenUtil *util.TokenUtil, mailer *util.Mailer) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
		Mailer:         mailer,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed count user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("User already exists : %+v", err)
		return nil, fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
		Email:    request.Email,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.TokenResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Invalid password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := c.TokenUtil.CreateAccessToken(&model.Auth{ID: user.ID})
	if err != nil {
		c.Log.Warnf("Failed to create token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.TokenUtil.CreateRefreshToken(&model.Auth{ID: user.ID})
	if err != nil {
		c.Log.Warnf("Failed to create refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.Token = refreshToken
	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed Commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	go func() {
		if err := c.Mailer.SendLoginNotification(user.Email); err != nil {
			c.Log.Warnf("Failed to send login notification to %s: %+v", user.Email, err)
		}
	}()

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *UserUseCase) Refresh(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	db := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	auth, err := c.TokenUtil.ParseToken(request.Token, "refresh")
	if err != nil {
		c.Log.Warnf("Invalid refresh token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	user := &entity.User{}
	if err := c.UserRepository.FindById(db, user, auth.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if user.Token != request.Token {
		c.Log.Warnf("Refresh token mismatch for user %s", auth.ID)
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := c.TokenUtil.CreateAccessToken(&model.Auth{ID: user.ID})
	if err != nil {
		c.Log.Warnf("Failed to create access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: user.Token,
	}, nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Faild find user by id %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed generate bcrypt hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed Commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}
