package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func VenueEntityToResponse(venue *entity.Venue) *model.VenueResponse {
	return &model.VenueResponse{
		ID:       venue.ID,
		Name:     venue.Name,
		Address:  venue.Address,
		Capacity: venue.Capacity,
		City:     venue.City,
		State:    venue.State,
		Zip:      venue.Zip,
	}
}

func VenuesToResponses(venues []entity.Venue) []*model.VenueResponse {
	venuesResponses := make([]*model.VenueResponse, len(venues))
	for i := range venues {
		venuesResponses[i] = VenueEntityToResponse(&venues[i])
	}
	return venuesResponses
}

func VenuesToPaginatedResponse(venues []entity.Venue, totalItems int64, page, size int) *model.Response[[]*model.VenueResponse] {
	venuesResponse := VenuesToResponses(venues)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(venuesResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
