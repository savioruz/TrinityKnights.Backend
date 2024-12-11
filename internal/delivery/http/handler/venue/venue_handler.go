package venue

import (
	"github.com/labstack/echo/v4"
)

type VenueHandler interface {
	CreateVenue(ctx echo.Context) error
	UpdateVenue(ctx echo.Context) error
	GetVenueByID(ctx echo.Context) error
	GetAllVenues(ctx echo.Context) error
	SearchVenues(ctx echo.Context) error
}
