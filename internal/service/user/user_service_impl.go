package user

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/user"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *user.UserRepositoryImpl
	JWTService     jwt.JWTService
	Viper		   *viper.Viper
	helper         *helper.ContextHelper
}

func NewUserServiceImpl(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *user.UserRepositoryImpl, jwtService jwt.JWTService, viper *viper.Viper) *UserServiceImpl {
	return &UserServiceImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTService:     jwtService,
		Viper:          viper,
		helper:         helper.NewContextHelper(),
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	existingUser := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, existingUser, request.Email); err == nil {
		s.Log.Errorf("email already exists: %v", request.Email)
		return nil, domainErrors.ErrDuplicateEntry
	}

	first := &entity.User{}
	if err := s.UserRepository.GetFirst(tx, first); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to get first user: %v", err)
		return nil, domainErrors.ErrInternalServer
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
		return nil, domainErrors.ErrInternalServer
	}

	

	token := uuid.NewString()
	var replaceEmail = struct {
		Token string
	}{
		Token: token,
	}

	

	
	data := &entity.User{
		ID:       uuid.NewString(),
		Email:    request.Email,
		Password: string(password),
		Name:     request.Name,
		Role:     role,
		VerifyEmailToken: token,
		IsVerify: false,
		Status:   true,
	}

	if err := s.UserRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	templatePath := "template/verify-email.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, &replaceEmail); err != nil {
		return nil, err
	}
	
	m := &model.SendEmail{
		EmailTo:   data.Email,
		EmailFrom: s.Viper.GetString("SMTP_USERNAME"),
		Subject:   "Reset Password Request",
		Body:      body,
	}
	if err := s.sendEmail(m); err != nil {
		return nil, err
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, data, request.Email); err != nil {
		s.Log.Errorf("failed to get user by email: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		s.Log.Errorf("failed to compare password: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := s.JWTService.GenerateAccessToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) Profile(ctx context.Context) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to get user by id: %v", err)
		return nil, domainErrors.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Update(ctx context.Context, request *model.UpdateRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil {
		s.Log.Errorf("failed to get user by id: %v", err)
		return nil, domainErrors.ErrNotFound
	}

	if request.Email != "" {
		data.Email = request.Email
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Log.Errorf("failed to generate password: %v", err)
			return nil, domainErrors.ErrInternalServer
		}
		data.Password = string(password)
	}

	if request.Name != "" {
		data.Name = request.Name
	}

	if err := s.UserRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	claims, err := s.JWTService.ValidateToken(request.RefreshToken)
	if err != nil {
		s.Log.Errorf("failed to validate token: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := s.JWTService.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) RequestReset(ctx context.Context, request *model.RequestReset) (*model.ResponseReset, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	existingUser := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, existingUser, request.Email); err != nil {
		s.Log.Errorf("email already exists: %v", request.Email)
		return nil, domainErrors.ErrNotFound
	}
	
	templatePath := "template/reset-password.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	token := uuid.NewString()
	var replaceEmail = struct {
		Token string
	}{
		Token: token,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, &replaceEmail); err != nil {
		return nil, err
	}

	m := &model.SendEmail{
		EmailTo:   existingUser.Email,
		EmailFrom: s.Viper.GetString("SMTP_USERNAME"),
		Subject:   "Reset Password Request",
		Body:      body,
	}
	if err := s.sendEmail(m); err != nil {
		return nil, err
	}

	update := &entity.User{
		ID:            existingUser.ID,
		Email:         existingUser.Email,
		Password:      existingUser.Password,
		Name:          existingUser.Name,
		Role:          existingUser.Role,
		Status:        existingUser.Status,
		ResetPasswordToken: token,
	}
	if err := s.UserRepository.Update(tx, update); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.ResponseReset{Status: "success"}, nil
}

func (s *UserServiceImpl) ResetPassword(ctx context.Context, request *model.ResetPassword) (*model.ResponseReset, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	user := &entity.User{}
	if err := s.UserRepository.GetByResetPasswordToken(tx, user, request.Token);err != nil {
		s.Log.Errorf("failed to get user by reset password token: %v", err)
		return nil, domainErrors.ErrNotFound
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("failed to generate password: %v", err)
		return nil, domainErrors.ErrInternalServer
	}
	
	user.Password = string(hashedPassword)
	err = s.UserRepository.Update(tx, user)
	if err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}
	

	return &model.ResponseReset{Status: "success"}, nil
}

func (s *UserServiceImpl) VerifyEmail(ctx context.Context, request *model.VerifyEmail) (*model.ResponseReset, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{}
	if err := s.UserRepository.GetByVerifyEmailToken(tx, user, request.Token);err != nil {
		s.Log.Errorf("failed to get user by reset password token: %v", err)
		return nil, domainErrors.ErrNotFound
	}

	if request.Token != user.VerifyEmailToken {
		return nil, domainErrors.ErrBadRequest
	}

	user.IsVerify = true
	if err := s.UserRepository.Update(tx, user); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.ResponseReset{Status: "success"}, nil
}

func (s *UserServiceImpl) sendEmail(request *model.SendEmail) error {
	m := gomail.NewMessage()
	m.SetHeader("From",request.EmailFrom)
	m.SetHeader("To",request.EmailTo)
	m.SetHeader("Subject", request.Subject)
	m.SetBody("text/html",request.Body.String())

	d := gomail.NewDialer(
		s.Viper.GetString("SMTP_HOST"),
		s.Viper.GetInt("SMTP_PORT"),
		s.Viper.GetString("SMTP_USERNAME"),
		s.Viper.GetString("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	
	return nil
}