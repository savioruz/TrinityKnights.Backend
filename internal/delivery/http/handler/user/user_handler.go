package user

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
	Profile(ctx echo.Context) error
	Update(ctx echo.Context) error
	RefreshToken(ctx echo.Context) error
	RequestReset(ctx echo.Context) error
	ResetPassword(ctx echo.Context) error
	VerifyEmail(ctx echo.Context) error
}
