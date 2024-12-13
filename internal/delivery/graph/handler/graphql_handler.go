package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/TrinityKnights/Backend/internal/delivery/graph"
	"github.com/TrinityKnights/Backend/internal/delivery/graph/resolvers"
	"github.com/labstack/echo/v4"
)

type GraphQLHandler struct {
	resolver *resolvers.Resolver
}

func NewGraphQLHandler(resolver *resolvers.Resolver) *GraphQLHandler {
	return &GraphQLHandler{
		resolver: resolver,
	}
}

// GraphQLHandler handles GraphQL requests
func (h *GraphQLHandler) GraphQLHandler(c echo.Context) error {
	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: h.resolver},
		),
	)

	graphqlHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}

// PlaygroundHandler serves the GraphQL playground interface
func (h *GraphQLHandler) PlaygroundHandler(c echo.Context) error {
	playgroundHandler := playground.Handler("GraphQL Playground", "/api/v1/graphql")
	playgroundHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}
