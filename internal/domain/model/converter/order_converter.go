package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

func OrderEntityToResponse(order *entity.Order) *model.OrderResponse {
	response := &model.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		Quantity:   len(order.Tickets),
	}

	if len(order.Tickets) > 0 {
		response.Tickets = make([]model.TicketResponse, len(order.Tickets))
		for i := range order.Tickets {
			ticket := &order.Tickets[i]
			response.Tickets[i] = model.TicketResponse{
				ID:         ticket.ID,
				EventID:    ticket.EventID,
				OrderID:    helper.UintOrZero(ticket.OrderID),
				Price:      ticket.Price,
				Type:       ticket.Type,
				SeatNumber: ticket.SeatNumber,
			}
			if i == 0 {
				response.EventID = ticket.EventID
			}
		}
	}

	return response
}

func OrdersToResponses(orders []entity.Order) []*model.OrderResponse {
	orderResponses := make([]*model.OrderResponse, len(orders))
	for i := range orders {
		orderResponses[i] = OrderEntityToResponse(&orders[i])
	}
	return orderResponses
}

func OrdersToPaginatedResponse(orders []entity.Order, totalItems int64, page, size int) *model.Response[[]*model.OrderResponse] {
	ordersResponse := OrdersToResponses(orders)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(ordersResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
