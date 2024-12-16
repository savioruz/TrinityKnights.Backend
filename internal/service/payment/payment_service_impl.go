package payment

import (
	"github.com/TrinityKnights/Backend/pkg/cache"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	DB       *gorm.DB
	Cache    *cache.ImplCache
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewPaymentServiceImpl(db *gorm.DB, cacheImpl *cache.ImplCache, log *logrus.Logger, validate *validator.Validate) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		DB:       db,
		Cache:    cacheImpl,
		Log:      log,
		Validate: validate,
	}
}
