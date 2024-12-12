package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func PaymentEntityToResponse(payment *entity.Payment) *model.PaymentResponse {
	return &model.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Method:        payment.Method,
		TransactionID: payment.TransactionID,
	}
}

func PaymentsToResponses(payments []entity.Payment) []*model.PaymentResponse {
	paymentsResponses := make([]*model.PaymentResponse, len(payments))
	for i := range payments {
		paymentsResponses[i] = PaymentEntityToResponse(&payments[i])
	}
	return paymentsResponses
}

func PaymentsToPaginatedResponse(payments []entity.Payment, totalItems int64, page, size int) *model.Response[[]*model.PaymentResponse] {
	paymentsResponse := PaymentsToResponses(payments)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(paymentsResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
