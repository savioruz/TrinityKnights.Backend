package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/TrinityKnights/Backend/internal/service/event"
	"github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/internal/service/venue"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

type Resolver struct {
	UserService  user.UserService
	EventService event.EventService
	VenueService venue.VenueService
	helper       helper.ContextHelper
}

func NewResolver(userService user.UserService, eventService event.EventService, venueService venue.VenueService) *Resolver {
	return &Resolver{
		UserService:  userService,
		EventService: eventService,
		VenueService: venueService,
		helper:       *helper.NewContextHelper(),
	}
}
