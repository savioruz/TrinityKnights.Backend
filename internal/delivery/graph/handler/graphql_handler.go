package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/TrinityKnights/Backend/internal/delivery/graph"
	"github.com/TrinityKnights/Backend/internal/delivery/graph/resolvers"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/labstack/echo/v4"
)

const contextKey = "claims"

type GraphQLHandler struct {
	resolver   *resolvers.Resolver
	jwtService jwt.JWTService
}

func NewGraphQLHandler(resolver *resolvers.Resolver, jwtService jwt.JWTService) *GraphQLHandler {
	return &GraphQLHandler{
		resolver:   resolver,
		jwtService: jwtService,
	}
}

func (h *GraphQLHandler) createGraphQLServer() *handler.Server {
	return handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: h.resolver,
				Directives: graph.DirectiveRoot{
					Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
						claims := ctx.Value(contextKey)
						if claims == nil {
							return nil, fmt.Errorf("access denied")
						}

						role := ctx.Value(contextKey).(*jwt.JWTClaims).Role
						if role != "admin" {
							return nil, fmt.Errorf("access denied")
						}

						return next(ctx)
					},
					Public: func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
						return next(ctx)
					},
				},
			},
		),
	)
}

// GraphQLHandler handles both public and private GraphQL requests
func (h *GraphQLHandler) GraphQLHandler(c echo.Context) error {
	// Check for Authorization header and validate token if present
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
			claims, err := h.jwtService.ValidateToken(bearerToken[1])
			if err == nil {
				// If token is valid, add claims to context
				ctx := context.WithValue(c.Request().Context(), contextKey, claims)
				c.SetRequest(c.Request().WithContext(ctx))
			}
		}
	}

	graphqlHandler := h.createGraphQLServer()
	graphqlHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}

// PlaygroundHandler serves the GraphQL playground interface
func (h *GraphQLHandler) PlaygroundHandler(c echo.Context) error {
	playgroundHandler := playground.Handler("GraphQL Playground", "/api/v1/graphql")
	playgroundHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}
