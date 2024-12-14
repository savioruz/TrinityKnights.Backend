package route

import (
	graphql "github.com/TrinityKnights/Backend/internal/delivery/graph/handler"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/event"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/order"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/payment"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/venue"
	"github.com/TrinityKnights/Backend/pkg/route"
	"github.com/labstack/echo/v4"
	swagger "github.com/swaggo/echo-swagger"
)

type RouteConfig interface {
	PublicRoute() []route.Route
	PrivateRoute() []route.Route
	SwaggerRoutes()
}

type Config struct {
	App            *echo.Echo
	GraphQLHandler *graphql.GraphQLHandler
	UserHandler    *user.UserHandlerImpl
	VenueHandler   *venue.VenueHandlerImpl
	EventHandler   *event.EventHandlerImpl
	OrderHandler   *order.OrderHandlerImpl
	PaymentHandler *payment.PaymentHandlerImpl
}

func (c Config) PublicRoute() []route.Route {
	return []route.Route{
		{
			Method:  echo.POST,
			Path:    "/users",
			Handler: c.UserHandler.Register,
		},
		{
			Method:  echo.POST,
			Path:    "/users/login",
			Handler: c.UserHandler.Login,
		},
		{
			Method:  echo.POST,
			Path:    "/users/refresh",
			Handler: c.UserHandler.RefreshToken,
		},
		{
			Method:  echo.POST,
			Path:    "/payments/callback",
			Handler: c.PaymentHandler.HandleCallback,
		},
	}
}

func (c Config) PrivateRoute() []route.Route {
	return []route.Route{
		{
			Method:  echo.GET,
			Path:    "/users",
			Handler: c.UserHandler.Profile,
		},
		{
			Method:  echo.PUT,
			Path:    "/users",
			Handler: c.UserHandler.Update,
		},
		{
			Method:  echo.POST,
			Path:    "/venues",
			Handler: c.VenueHandler.CreateVenue,
		},
		{
			Method:  echo.GET,
			Path:    "/venues",
			Handler: c.VenueHandler.GetAllVenues,
		},
		{
			Method:  echo.GET,
			Path:    "/venues/:id",
			Handler: c.VenueHandler.GetVenueByID,
		},
		{
			Method:  echo.PUT,
			Path:    "/venues/:id",
			Handler: c.VenueHandler.UpdateVenue,
		},
		{
			Method:  echo.GET,
			Path:    "/venues/search",
			Handler: c.VenueHandler.SearchVenues,
		},
		{
			Method:  echo.POST,
			Path:    "/events",
			Handler: c.EventHandler.CreateEvent,
		},
		{
			Method:  echo.PUT,
			Path:    "/events/:id",
			Handler: c.EventHandler.UpdateEvent,
		},
		{
			Method:  echo.GET,
			Path:    "/events/:id",
			Handler: c.EventHandler.GetEventByID,
		},
		{
			Method:  echo.GET,
			Path:    "/events",
			Handler: c.EventHandler.GetAllEvents,
		},
		{
			Method:  echo.GET,
			Path:    "/events/search",
			Handler: c.EventHandler.SearchEvents,
		},
		{
			Method:  echo.POST,
			Path:    "/orders",
			Handler: c.OrderHandler.CreateOrder,
		},
		{
			Method:  echo.GET,
			Path:    "/orders/:id",
			Handler: c.OrderHandler.GetOrderByID,
		},
		{
			Method:  echo.GET,
			Path:    "/orders",
			Handler: c.OrderHandler.GetAllOrders,
		},
		{
			Method:  echo.POST,
			Path:    "/payments",
			Handler: c.PaymentHandler.CreatePayment,
		},
	}
}

func (c Config) SwaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}
