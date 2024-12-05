package main

import (
	"time"

	"github.com/TrinityKnights/Backend/config"
	"github.com/TrinityKnights/Backend/internal/builder"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
func main() {

	viper := config.NewViper()

	log := config.NewLogrus()

	jwtService := jwt.NewJWTService(config.NewJWT(viper))

	db := config.NewDatabase(viper, log)

	redisClient := config.NewRedisClient(viper, log)

	e, log := config.NewEcho()
	validator := config.NewValidator()
	e.Validator = &CustomValidator{validator: validator}

	routes := builder.BuildPublicRoutes(log, db, jwtService)

	for _, route := range routes {
		e.Add(route.Method, route.Path, route.Handler)
	}

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	log.Println("Server started on port 8080")

	config.GracefulShutdown(e, log, 10*time.Second)

	log.Info("Server is running and connected to Redis")
	redisClient.Set("exampleKey", "exampleValue", 0)
	err := redisClient.Get("exampleKey", nil)
	if err != nil {
		log.Errorf("Error getting key from Redis: %v", err)
	} else {
		log.Info("Successfully got key from Redis")
	}
}
