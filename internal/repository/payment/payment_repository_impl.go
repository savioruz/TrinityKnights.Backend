package payment

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewPaymentRepository(db *gorm.DB, log *logrus.Logger) PaymentRepository {
	return &PaymentRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *PaymentRepositoryImpl) GetByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePaymentStatus(ctx context.Context, payment *model.PaymentUpdateRequest) error {
	return r.DB.Model(&entity.Payment{}).Where("id = ?", payment.ID).Updates(payment).Error
}

func (r *PaymentRepositoryImpl) Find(db *gorm.DB, filter *model.PaymentQueryOptions) ([]*entity.Payment, error) {
	query := db.Model(&entity.Payment{})

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}

	if filter.Method != nil {
		query = query.Where("UPPER(method) = UPPER(?)", *filter.Method)
	}

	if filter.TransactionID != nil {
		query = query.Where("transaction_id = ?", *filter.TransactionID)
	}

	if filter.Status != nil {
		query = query.Where("UPPER(status) = UPPER(?)", *filter.Status)
	}

	var payments []*entity.Payment
	if err := query.Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
