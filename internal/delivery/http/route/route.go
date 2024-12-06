package route

import (
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/pkg/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	App            *echo.Echo
	UserHandler    *user.UserHandlerImpl
	AuthMiddleware echo.MiddlewareFunc
}

func (c *Config) Setup() {
	g := c.App.Group("/api/v1")
	for _, r := range c.publicRoute() {
		g.Add(r.Method, r.Path, r.Handler)
	}
	for _, r := range c.privateRoute() {
		g.Add(r.Method, r.Path, r.Handler, c.AuthMiddleware)
	}
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(60)))
	c.swaggerRoutes()
	c.App.Use(middleware.Recover())
}

func (c *Config) publicRoute() []route.Route {
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
	}
}

func (c *Config) privateRoute() []route.Route {
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
	}
}

func (c *Config) swaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}
