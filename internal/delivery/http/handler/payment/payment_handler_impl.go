package payment

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
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

// @Summary Create a new payment
// @Description Create a new payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param request body model.PaymentRequest true "Payment details"
// @Success 201 {object} model.Response[model.PaymentResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /payments [post]
func (h *PaymentHandlerImpl) CreatePayment(ctx echo.Context) error {
	request := new(model.PaymentRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.PaymentService.CreatePayment(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create payment: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// @Summary Handle payment callback
// @Description Handle payment callback from payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param request body model.MidtransCallbackRequest true "Callback details"
// @Success 200 {object} model.Response[any]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /payments/callback [post]
func (h *PaymentHandlerImpl) HandleCallback(ctx echo.Context) error {
	request := new(model.MidtransCallbackRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind callback request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	if err := h.PaymentService.HandleCallback(ctx.Request().Context(), request); err != nil {
		h.Log.Errorf("failed to handle callback: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse[any](nil, nil))
}
