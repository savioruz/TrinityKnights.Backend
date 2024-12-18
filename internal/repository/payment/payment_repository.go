package payment

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type PaymentRepository interface {
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, payment *model.PaymentUpdateRequest) error
}
