package payment

import (
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	"github.com/labstack/echo/v4"
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

// @Summary Callback Payment
// @Description Callback Payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param payment body model.PaymentCallbackRequest true "Payment Callback Request"
// @Success 200 {object} model.PaymentCallbackResponse
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /payment/callback [post]
func (h *PaymentHandlerImpl) CallbackPayment(ctx echo.Context) error {
	request := new(model.PaymentCallbackRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.PaymentService.Callback(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to callback payment: %v", err)
		return handler.HandleError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
