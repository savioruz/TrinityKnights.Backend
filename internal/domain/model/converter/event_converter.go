package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func EventEntityToResponse(event *entity.Event) *model.EventResponse {
	return &model.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Time:        event.Time,
		VenueID:     event.VenueID,
	}
}

func EventsToResponses(events []entity.Event) []*model.EventResponse {
	eventResponses := make([]*model.EventResponse, len(events))
	for i := range events {
		eventResponses[i] = EventEntityToResponse(&events[i])
	}
	return eventResponses
}

func EventsToPaginatedResponse(events []entity.Event, totalItems int64, page, size int) *model.Response[[]*model.EventResponse] {
	eventsResponse := EventsToResponses(events)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(eventsResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
