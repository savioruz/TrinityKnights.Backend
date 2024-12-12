package payment

import (
	"context"
	"fmt"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	DB       *gorm.DB
	Cache    *cache.ImplCache
	Log      *logrus.Logger
	Validate *validator.Validate
	Snap     snap.Client
}

func NewPaymentServiceImpl(db *gorm.DB, cache *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, snap snap.Client) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		DB:       db,
		Cache:    cache,
		Log:      log,
		Validate: validate,
		Snap:     snap,
	}
}

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, request *model.PaymentRequest) (*model.PaymentResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Create Midtrans transaction
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("ORDER-%d", request.OrderID),
			GrossAmt: int64(request.Amount),
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeBankTransfer,
		},
	}

	snapResp, err := s.Snap.CreateTransaction(snapReq)
	if err != nil {
		s.Log.Errorf("failed to create midtrans transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create payment record
	payment := &entity.Payment{
		OrderID:       request.OrderID,
		Method:        request.PaymentMethod,
		TransactionID: snapResp.Token,
	}

	if err := tx.Create(payment).Error; err != nil {
		s.Log.Errorf("failed to create payment: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return &model.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        request.Amount,
		Status:        "PENDING",
		PaymentMethod: payment.Method,
		MidtransID:    payment.TransactionID,
		PaymentURL:    snapResp.RedirectURL,
	}, nil
}

func (s *PaymentServiceImpl) HandleCallback(ctx context.Context, request *model.MidtransCallbackRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var payment entity.Payment
	if err := tx.Where("transaction_id = ?", request.OrderID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get payment: %v", err)
		return domainErrors.ErrInternalServer
	}

	// Update payment status based on callback
	// Implementation depends on your business logic

	if err := tx.Commit().Error; err != nil {
		return domainErrors.ErrInternalServer
	}

	return nil
}
