package payment

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, request *model.PaymentRequest) (*model.PaymentResponse, error)
	HandleCallback(ctx context.Context, request *model.MidtransCallbackRequest) error
}
