package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/user"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *user.UserRepositoryImpl
	JWTService     jwt.JWTService
	helper         *helper.ContextHelper
}

func NewUserServiceImpl(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *user.UserRepositoryImpl, jwtService jwt.JWTService) *UserServiceImpl {
	return &UserServiceImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTService:     jwtService,
		helper:         helper.NewContextHelper(),
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	existingUser := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, existingUser, request.Email); err == nil {
		s.Log.Errorf("email already exists: %v", request.Email)
		return nil, errors.New(http.StatusText(http.StatusConflict))
	}

	first := &entity.User{}
	if err := s.UserRepository.GetFirst(tx, first); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to get first user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	var role string
	if first.ID == "" {
		role = "admin"
	} else {
		role = "buyer"
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("failed to generate password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	data := &entity.User{
		ID:       uuid.New().String(),
		Email:    request.Email,
		Password: string(password),
		Name:     request.Name,
		Role:     role,
		Status:   true,
	}

	if err := s.UserRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, data, request.Email); err != nil {
		s.Log.Errorf("failed to get user by email: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		s.Log.Errorf("failed to compare password: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	accessToken, err := s.JWTService.GenerateAccessToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) Profile(ctx context.Context) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to get user by id: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Update(ctx context.Context, request *model.UpdateRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil {
		s.Log.Errorf("failed to get user by id: %v", err)
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if request.Email != "" {
		data.Email = request.Email
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Log.Errorf("failed to generate password: %v", err)
			return nil, errors.New(http.StatusText(http.StatusInternalServerError))
		}
		data.Password = string(password)
	}

	if request.Name != "" {
		data.Name = request.Name
	}

	if err := s.UserRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	claims, err := s.JWTService.ValidateToken(request.RefreshToken)
	if err != nil {
		s.Log.Errorf("failed to validate token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	accessToken, err := s.JWTService.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) ResetPasswordRequest(ctx context.Context, request *model.ResetPasswordRequest) error {
	return nil
}