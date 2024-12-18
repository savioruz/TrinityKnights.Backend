package payment

import (
	"github.com/labstack/echo/v4"
)

type PaymentHandler interface {
	CallbackPayment(ctx echo.Context) error
}
