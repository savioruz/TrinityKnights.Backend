package handler

import (
	"net/http"

	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) UserHandler {
	return UserHandler{userService}
}

func (h UserHandler) Register(ctx echo.Context) error {
	request := new(model.RegisterRequest)
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := ctx.Validate(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	userResponse, err := h.userService.Register(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusCreated, response.SuccessResponse("user created successfully", map[string]interface{}{
		"user": userResponse,
	}))
}

func (h UserHandler) Login(ctx echo.Context) error {
	request := new(model.LoginRequest)
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if err := ctx.Validate(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	tokenResponse, err := h.userService.Login(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("login successful", map[string]interface{}{
		"token": tokenResponse,
	}))
}
