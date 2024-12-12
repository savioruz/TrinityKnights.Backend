package event

import (
	"github.com/labstack/echo/v4"
)

type EventHandler interface {
	CreateEvent(ctx echo.Context) error
	UpdateEvent(ctx echo.Context) error
	GetEventByID(ctx echo.Context) error
	GetAllEvents(ctx echo.Context) error
	SearchEvents(ctx echo.Context) error
}
