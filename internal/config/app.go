package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rendi-hendra/resful-api/internal/delivery/http"
	"github.com/rendi-hendra/resful-api/internal/delivery/http/middleware"
	"github.com/rendi-hendra/resful-api/internal/delivery/http/route"
	"github.com/rendi-hendra/resful-api/internal/repository"
	"github.com/rendi-hendra/resful-api/internal/usecase"
	"github.com/rendi-hendra/resful-api/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstapConfig) {
	// setup repository
	userRepository := repository.NewUserRepository(config.Log)

	// token
	tokenUtil := util.NewTokenUtil("rahasia")

	// setup usecase
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenUtil)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()
}
