package ticket

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/ticket"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TicketHandlerImpl struct {
	Log           *logrus.Logger
	TicketService ticket.TicketService
}

func NewTicketHandler(log *logrus.Logger, ticketService ticket.TicketService) TicketHandler {
	return &TicketHandlerImpl{
		Log:           log,
		TicketService: ticketService,
	}
}

// @Summary Create new tickets @admin
// @Description Create new tickets with the provided details
// @Tags tickets
// @Accept json
// @Produce json
// @Param request body model.CreateTicketRequest true "Ticket details"
// @Success 201 {object} model.Response[[]model.TicketResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /tickets [post]
func (h *TicketHandlerImpl) CreateTicket(ctx echo.Context) error {
	request := new(model.CreateTicketRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.TicketService.CreateTicket(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create tickets: %v", err)
		if errors.Is(err, domainErrors.ErrBadRequest) {
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		}
		return handler.HandleError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// @Summary Update an existing ticket @admin
// @Description Update an existing ticket with the provided details
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param request body model.UpdateTicketRequest true "Updated ticket details"
// @Success 200 {object} model.Response[model.TicketResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /tickets/{id} [put]
func (h *TicketHandlerImpl) UpdateTicket(ctx echo.Context) error {
	request := new(model.UpdateTicketRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.TicketService.UpdateTicket(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to update ticket: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get a ticket by ID
// @Description Get details of a specific ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} model.Response[model.TicketResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /tickets/{id} [get]
func (h *TicketHandlerImpl) GetTicketByID(ctx echo.Context) error {
	request := new(model.GetTicketRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.TicketService.GetTicketByID(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get ticket by id: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get all tickets
// @Description Get a paginated list of all tickets
// @Tags tickets
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(event_id,order_id,price,type,seat_number)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.TicketResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /tickets [get]
func (h *TicketHandlerImpl) GetAllTickets(ctx echo.Context) error {
	request := new(model.TicketsRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.TicketService.GetTickets(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get tickets: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// @Summary Search tickets
// @Description Search tickets with the provided query parameters
// @Tags tickets
// @Produce json
// @Param id query string false "Ticket ID"
// @Param event_id query int false "Event ID"
// @Param order_id query int false "Order ID"
// @Param price query number false "Ticket price"
// @Param type query string false "Ticket type"
// @Param seat_number query string false "Seat number"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(id, event_id, order_id, price, type, seat_number)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.TicketResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /tickets/search [get]
func (h *TicketHandlerImpl) SearchTickets(ctx echo.Context) error {
	request := new(model.TicketSearchRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.TicketService.SearchTickets(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to search tickets: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
