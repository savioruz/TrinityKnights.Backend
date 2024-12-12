package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

// OrderEntityToResponse mengonversi entitas Order menjadi OrderResponse.
func OrderEntityToResponse(order *entity.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		Date:       order.Date,
		TotalPrice: order.TotalPrice,
		PaymentID:  order.PaymentID,
		// Konversi Tickets
	}
}

// OrdersToResponses mengonversi daftar Order menjadi slice OrderResponse.
func OrdersToResponses(orders []entity.Order) []*model.OrderResponse {
	orderResponses := make([]*model.OrderResponse, len(orders))
	for i := range orders {
		orderResponses[i] = OrderEntityToResponse(&orders[i])
	}
	return orderResponses
}

// OrdersToPaginatedResponse mengonversi daftar Order menjadi respons dengan pagination.
func OrdersToPaginatedResponse(orders []entity.Order, totalItems int64, page, size int) *model.Response[[]*model.OrderResponse] {
	orderResponses := OrdersToResponses(orders)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(orderResponses, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
