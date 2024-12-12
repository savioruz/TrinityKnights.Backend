package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func OrderEntityToResponse(order *entity.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		PaymentID:  order.PaymentID,
	}
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
