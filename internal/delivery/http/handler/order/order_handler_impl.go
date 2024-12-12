package order

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/order"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type OrderHandlerImpl struct {
	Log          *logrus.Logger
	OrderService order.OrderService
}

func NewOrderHandler(log *logrus.Logger, orderService order.OrderService) OrderHandler {
	return &OrderHandlerImpl{
		Log:          log,
		OrderService: orderService,
	}
}

// @Summary Create a new order
// @Description Create a new order for event tickets
// @Tags orders
// @Accept json
// @Produce json
// @Param request body model.OrderTicketRequest true "Order details"
// @Success 201 {object} model.Response[model.OrderResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /orders [post]
func (h *OrderHandlerImpl) CreateOrder(ctx echo.Context) error {
	request := new(model.OrderTicketRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.OrderService.CreateOrder(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create order: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// @Summary Get order by ID
// @Description Get details of a specific order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.Response[model.OrderResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /orders/{id} [get]
func (h *OrderHandlerImpl) GetOrderByID(ctx echo.Context) error {
	request := new(model.GetOrderRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.OrderService.GetOrderByID(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get order: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get all orders
// @Description Get a paginated list of all orders
// @Tags orders
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(date, total_price)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.OrderResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /orders [get]
func (h *OrderHandlerImpl) GetAllOrders(ctx echo.Context) error {
	request := new(model.OrdersRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.OrderService.GetOrders(ctx.Request().Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
