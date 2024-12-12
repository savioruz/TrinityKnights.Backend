package payment

import (
	"github.com/labstack/echo/v4"
)

type PaymentHandler interface {
	CreatePayment(ctx echo.Context) error
	HandleCallback(ctx echo.Context) error
}
