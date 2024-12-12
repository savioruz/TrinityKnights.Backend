package event

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type EventService interface {
	CreateEvent(ctx context.Context, request *model.CreateEventRequest) (*model.EventResponse, error)
	UpdateEvent(ctx context.Context, request *model.UpdateEventRequest) (*model.EventResponse, error)
	GetEventByID(ctx context.Context, request *model.GetEventRequest) (*model.EventResponse, error)
	GetEvents(ctx context.Context, request *model.EventsRequest) (*model.Response[[]*model.EventResponse], error)
	SearchEvents(ctx context.Context, request *model.EventSearchRequest) (*model.Response[[]*model.EventResponse], error)
}
