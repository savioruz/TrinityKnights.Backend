package ticket

import (
	"github.com/labstack/echo/v4"
)

type TicketHandler interface {
	CreateTicket(ctx echo.Context) error
	UpdateTicket(ctx echo.Context) error
	GetTicketByID(ctx echo.Context) error
	GetAllTickets(ctx echo.Context) error
	SearchTickets(ctx echo.Context) error
}
