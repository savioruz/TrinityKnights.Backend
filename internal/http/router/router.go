package router

import (
	"net/http"

	"github.com/TrinityKnights/Backend/internal/http/handler"
	"github.com/TrinityKnights/Backend/pkg/route"
)

func PublicRoutes(userHandler handler.UserHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
	}
}
