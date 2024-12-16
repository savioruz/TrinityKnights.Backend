package ticket

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type TicketService interface {
	CreateTicket(ctx context.Context, request *model.CreateTicketRequest) ([]*model.TicketResponse, error)
	UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) (*model.TicketResponse, error)
	GetTicketByID(ctx context.Context, request *model.GetTicketRequest) (*model.TicketResponse, error)
	GetTickets(ctx context.Context, request *model.TicketsRequest) (*model.Response[[]*model.TicketResponse], error)
	SearchTickets(ctx context.Context, request *model.TicketSearchRequest) (*model.Response[[]*model.TicketResponse], error)
}
