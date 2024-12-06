package config

import (
	handlerUser "github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/delivery/http/middleware"
	"github.com/TrinityKnights/Backend/internal/delivery/http/route"
	repositoryUser "github.com/TrinityKnights/Backend/internal/repository/user"
	serviceUser "github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/pkg/cache"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Cache    *cache.ImplCache
	App      *echo.Echo
	Log      *logrus.Logger
	Validate *validator.Validate
	JWT      *jwt.JWTConfig
}

func Bootstrap(config *BootstrapConfig) error {
	// Initialize JWT service
	jwtService := jwt.NewJWTService(config.JWT)

	// Initialize repository
	userRepository := repositoryUser.NewUserRepository(config.DB, config.Log)

	// Initialize service
	userService := serviceUser.NewUserServiceImpl(config.DB, config.Log, config.Validate, userRepository, jwtService)

	// Initialize handler
	userHandler := handlerUser.NewUserHandler(config.Log, userService)

	// Initialize graphql

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Initialize route
	routeConfig := &route.Config{
		App:            config.App,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()

	config.Log.Infof("Application is ready")
	return nil
}
