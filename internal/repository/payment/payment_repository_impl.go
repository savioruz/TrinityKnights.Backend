package payment

import (
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
