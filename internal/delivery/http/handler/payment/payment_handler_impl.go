package payment

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaymentHandlerImpl struct {
	Viper          *viper.Viper
	Log            *logrus.Logger
	PaymentService payment.PaymentService
}

func NewPaymentHandler(viper *viper.Viper, log *logrus.Logger, paymentService payment.PaymentService) PaymentHandler {
	return &PaymentHandlerImpl{
		Viper:          viper,
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

	myCallbackToken := h.Viper.GetString("XENDIT_CALLBACK_TOKEN")

	// Verify
	callbackToken := ctx.Request().Header.Get("x-callback-token")

	if callbackToken != myCallbackToken {
		h.Log.Errorf("invalid webhook id or callback token")
		return handler.HandleError(ctx, http.StatusUnauthorized, errors.New("invalid webhook id or callback token"))
	}

	response, err := h.PaymentService.Callback(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to callback payment: %v", err)
		return handler.HandleError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
