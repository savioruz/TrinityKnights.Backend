package user

import (
	"context"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *user.UserRepositoryImpl
	// @TODO: JWT
}

func NewUserServiceImpl(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *user.UserRepositoryImpl) *UserServiceImpl {
	return &UserServiceImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	return nil, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	return nil, nil
}

func (s *UserServiceImpl) Profile(ctx context.Context, request *model.ProfileRequest) (*model.UserResponse, error) {
	return nil, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, request *model.UpdateRequest) (*model.UserResponse, error) {
	return nil, nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	return nil, nil
}
