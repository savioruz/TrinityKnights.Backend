package payment

import (
	"github.com/TrinityKnights/Backend/internal/service/payment"
	"github.com/sirupsen/logrus"
)

type PaymentHandlerImpl struct {
	Log            *logrus.Logger
	PaymentService payment.PaymentService
}

func NewPaymentHandler(log *logrus.Logger, paymentService payment.PaymentService) PaymentHandler {
	return &PaymentHandlerImpl{
		Log:            log,
		PaymentService: paymentService,
	}
}
