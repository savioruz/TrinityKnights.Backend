package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

func TicketEntityToResponse(ticket *entity.Ticket) *model.TicketResponse {
	return &model.TicketResponse{
		ID:         ticket.ID,
		EventID:    ticket.EventID,
		OrderID:    helper.UintOrZero(ticket.OrderID),
		Price:      ticket.Price,
		Type:       ticket.Type,
		SeatNumber: ticket.SeatNumber,
		Event:      EventEntityToResponse(&ticket.Event),
		Order:      OrderEntityToResponse(&ticket.Order),
	}
}

func TicketsToResponses(tickets []*entity.Ticket) []*model.TicketResponse {
	ticketsResponses := make([]*model.TicketResponse, len(tickets))
	for i := range tickets {
		ticketsResponses[i] = TicketEntityToResponse(tickets[i])
	}
	return ticketsResponses
}

func TicketsToPaginatedResponse(tickets []*entity.Ticket, totalItems int64, page, size int) *model.Response[[]*model.TicketResponse] {
	ticketsResponse := TicketsToResponses(tickets)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(ticketsResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
