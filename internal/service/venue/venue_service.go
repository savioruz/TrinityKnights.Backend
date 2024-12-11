package venue

import (
	"context"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type VenueService interface {
	CreateVenue(ctx context.Context, request *model.CreateVenueRequest) (*model.VenueResponse, error)
	UpdateVenue(ctx context.Context, request *model.UpdateVenueRequest) (*model.VenueResponse, error)
	GetVenueByID(ctx context.Context, request *model.GetVenueRequest) (*model.VenueResponse, error)
	GetVenues(ctx context.Context, request *model.VenuesRequest) (*model.Response[[]*model.VenueResponse], error)
	SearchVenues(ctx context.Context, request *model.VenueSearchRequest) (*model.Response[[]*model.VenueResponse], error)
}
