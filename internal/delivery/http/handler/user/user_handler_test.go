package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	mockUser "github.com/TrinityKnights/Backend/test/mock/service/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*user.UserHandlerImpl, *mockUser.MockUserService, *echo.Echo) {
	ctrl := gomock.NewController(t)
	mockUserService := mockUser.NewMockUserService(ctrl)
	logger := logrus.New()
	handler := user.NewUserHandler(logger, mockUserService)
	e := echo.New()
	return handler, mockUserService, e
}

func TestUserHandler_Register(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"email": "test@example.com",
				"password": "password123",
				"name": "Test User"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Register(gomock.Any(), &model.RegisterRequest{
						Email:    "test@example.com",
						Password: "password123",
						Name:     "Test User",
					}).
					Return(&model.UserResponse{
						ID:        "1",
						Email:     "test@example.com",
						Name:      "Test User",
						Role:      "",
						Status:    false,
						CreatedAt: "",
						UpdatedAt: "",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"data":{"id":"1","email":"test@example.com","name":"Test User","role":"","status":false,"created_at":"","updated_at":""}}`,
		},
		{
			name: "Email Already Exists",
			requestBody: `{
				"email": "existing@example.com",
				"password": "password123",
				"name": "Test User"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrEmailAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":{"code":409, "message":"email already exists"}}`,
		},
		{
			name:           "Invalid JSON Request",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "Empty Request",
			requestBody: `{
				"email": "",
				"password": "",
				"name": ""
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.Register(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"email": "test@example.com",
				"password": "password123"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), &model.LoginRequest{
						Email:    "test@example.com",
						Password: "password123",
					}).
					Return(&model.TokenResponse{
						AccessToken:  "access_token",
						RefreshToken: "refresh_token",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"access_token":"access_token","refresh_token":"refresh_token"}}`,
		},
		{
			name: "Invalid Credentials",
			requestBody: `{
				"email": "test@example.com",
				"password": "wrongpassword"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":401, "message":"unauthorized"}}`,
		},
		{
			name:           "Invalid JSON Request",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "Empty Request",
			requestBody: `{
				"email": "",
				"password": ""
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), &model.LoginRequest{
						Email:    "",
						Password: "",
					}).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.Login(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_Profile(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockUserService.EXPECT().
					Profile(gomock.Any()).
					Return(&model.UserResponse{
						ID:        "1",
						Email:     "test@example.com",
						Name:      "Test User",
						Role:      "",
						Status:    false,
						CreatedAt: "",
						UpdatedAt: "",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":"1","email":"test@example.com","name":"Test User","role":"","status":false,"created_at":"","updated_at":""}}`,
		},
		{
			name: "User Not Found",
			setupMock: func() {
				mockUserService.EXPECT().
					Profile(gomock.Any()).
					Return(nil, domainErrors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":404, "message":"not found"}}`,
		},
		{
			name: "Internal Server Error",
			setupMock: func() {
				mockUserService.EXPECT().
					Profile(gomock.Any()).
					Return(nil, domainErrors.ErrInternalServer)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":500,"message":"internal server error"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.Profile(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"name": "Updated Name",
				"email": "updated@example.com"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Update(gomock.Any(), &model.UpdateUserRequest{
						Name:  "Updated Name",
						Email: "updated@example.com",
					}).
					Return(&model.UserResponse{
						ID:        "1",
						Email:     "updated@example.com",
						Name:      "Updated Name",
						Role:      "user",
						Status:    true,
						CreatedAt: "2024-01-01T00:00:00Z",
						UpdatedAt: "2024-01-01T00:00:00Z",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":"1","email":"updated@example.com","name":"Updated Name","role":"user","status":true,"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}}`,
		},
		{
			name: "Invalid Request",
			requestBody: `{
				"name": "",
				"email": "invalid-email"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "User Not Found",
			requestBody: `{
				"name": "Updated Name",
				"email": "updated@example.com"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":404,"message":"not found"}}`,
		},
		{
			name: "Internal Server Error",
			requestBody: `{
				"name": "Updated Name",
				"email": "updated@example.com"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrInternalServer)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":{"code":500,"message":"internal server error"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.Update(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_RefreshToken(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"refresh_token": "valid_refresh_token"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RefreshToken(gomock.Any(), &model.RefreshTokenRequest{
						RefreshToken: "valid_refresh_token",
					}).
					Return(&model.TokenResponse{
						AccessToken:  "new_access_token",
						RefreshToken: "new_refresh_token",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"access_token":"new_access_token","refresh_token":"new_refresh_token"}}`,
		},
		{
			name: "Invalid Token",
			requestBody: `{
				"refresh_token": "invalid_token"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RefreshToken(gomock.Any(), &model.RefreshTokenRequest{
						RefreshToken: "invalid_token",
					}).
					Return(nil, domainErrors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":401,"message":"unauthorized"}}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "Empty Token",
			requestBody: `{
				"refresh_token": ""
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RefreshToken(gomock.Any(), &model.RefreshTokenRequest{
						RefreshToken: "",
					}).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/refresh", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.RefreshToken(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_RequestReset(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"email": "test@example.com"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RequestReset(gomock.Any(), &model.ReqResetPasswordRequest{
						Email: "test@example.com",
					}).
					Return(&model.VerifyResponse{
						Status: "reset password request sent successfully",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"status":"reset password request sent successfully"}}`,
		},
		{
			name: "User Not Found",
			requestBody: `{
				"email": "nonexistent@example.com"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RequestReset(gomock.Any(), &model.ReqResetPasswordRequest{
						Email: "nonexistent@example.com",
					}).
					Return(nil, domainErrors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":404,"message":"not found"}}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "Empty Email",
			requestBody: `{
				"email": ""
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					RequestReset(gomock.Any(), &model.ReqResetPasswordRequest{
						Email: "",
					}).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/reset-password/request", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.RequestReset(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_ResetPassword(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"token": "valid_reset_token",
				"email": "test@example.com",
				"new_password": "newpassword123"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					ResetPassword(gomock.Any(), &model.ResetPasswordRequest{
						Token:       "valid_reset_token",
						Email:       "test@example.com",
						NewPassword: "newpassword123",
					}).
					Return(&model.VerifyResponse{
						Status: "success",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"status":"success"}}`,
		},
		{
			name: "Invalid Token",
			requestBody: `{
				"token": "invalid_token",
				"email": "test@example.com",
				"new_password": "newpassword123"
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					ResetPassword(gomock.Any(), &model.ResetPasswordRequest{
						Token:       "invalid_token",
						Email:       "test@example.com",
						NewPassword: "newpassword123",
					}).
					Return(nil, domainErrors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":401,"message":"unauthorized"}}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
		{
			name: "Empty Fields",
			requestBody: `{
				"token": "",
				"password": ""
			}`,
			setupMock: func() {
				mockUserService.EXPECT().
					ResetPassword(gomock.Any(), &model.ResetPasswordRequest{
						Token:       "",
						NewPassword: "",
					}).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/reset-password", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.ResetPassword(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestUserHandler_VerifyEmail(t *testing.T) {
	handler, mockUserService, e := setupTest(t)

	tests := []struct {
		name           string
		token          string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "Success",
			token: "valid_verification_token",
			setupMock: func() {
				mockUserService.EXPECT().
					VerifyEmail(gomock.Any(), &model.VerifyRequest{
						Token: "valid_verification_token",
					}).
					Return(&model.VerifyResponse{
						Status: "success",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"status":"success"}}`,
		},
		{
			name:  "Invalid Token",
			token: "invalid_token",
			setupMock: func() {
				mockUserService.EXPECT().
					VerifyEmail(gomock.Any(), &model.VerifyRequest{
						Token: "invalid_token",
					}).
					Return(nil, domainErrors.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":{"code":401,"message":"unauthorized"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/verify-email/"+tc.token, http.NoBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set path parameter
			c.SetParamNames("token")
			c.SetParamValues(tc.token)

			tc.setupMock()

			err := handler.VerifyEmail(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}
