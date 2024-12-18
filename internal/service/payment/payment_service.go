package payment

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreateInvoice(ctx context.Context, tx *gorm.DB, request *model.PaymentRequest) (*model.PaymentResponse, error)
	Callback(ctx context.Context, request *model.PaymentCallbackRequest) (*model.PaymentCallbackResponse, error)
}
