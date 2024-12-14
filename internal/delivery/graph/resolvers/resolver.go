package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

type Resolver struct {
	UserService user.UserService
	helper      helper.ContextHelper
}

func NewResolver(userService user.UserService) *Resolver {
	return &Resolver{
		UserService: userService,
		helper:      *helper.NewContextHelper(),
	}
}
