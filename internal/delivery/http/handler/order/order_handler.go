package order

import (
	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	CreateOrder(ctx echo.Context) error
	GetOrderByID(ctx echo.Context) error
	GetAllOrders(ctx echo.Context) error
}
