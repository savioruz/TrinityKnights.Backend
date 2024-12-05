package builder

import (
	"github.com/TrinityKnights/Backend/internal/http/handler"
	"github.com/TrinityKnights/Backend/internal/http/router"
	userRepo "github.com/TrinityKnights/Backend/internal/repository/user"
	userService "github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/TrinityKnights/Backend/pkg/route"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func BuildPublicRoutes(log *logrus.Logger, db *gorm.DB, jwtService jwt.JWTService) []route.Route {
	validate := validator.New()
	userRepository := userRepo.NewUserRepository(db, log)
	userService := userService.NewUserServiceImpl(db, log, validate, userRepository, jwtService)
	userHandler := handler.NewUserHandler(userService)

	return router.PublicRoutes(userHandler)
}
