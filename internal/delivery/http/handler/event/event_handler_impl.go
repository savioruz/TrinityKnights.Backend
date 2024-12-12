package event

import (
	"errors"
	"net/http"
	"strings"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/event"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type EventHandlerImpl struct {
	Log          *logrus.Logger
	EventService event.EventService
}

func NewEventHandler(log *logrus.Logger, eventService event.EventService) EventHandler {
	return &EventHandlerImpl{
		Log:          log,
		EventService: eventService,
	}
}

// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param request body model.CreateEventRequest true "Event details"
// @Success 201 {object} model.Response[model.EventResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /events [post]
func (h *EventHandlerImpl) CreateEvent(ctx echo.Context) error {
	request := new(model.CreateEventRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.EventService.CreateEvent(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create event: %v", err)
		if strings.Contains(err.Error(), "invalid request") {
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		}
		return handler.HandleError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// @Summary Update an existing event
// @Description Update an existing event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param request body model.UpdateEventRequest true "Updated event details"
// @Success 200 {object} model.Response[model.EventResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /events/{id} [put]
func (h *EventHandlerImpl) UpdateEvent(ctx echo.Context) error {
	request := new(model.UpdateEventRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.EventService.UpdateEvent(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to update event: %v", err)
		switch {
		case errors.Is(err, errors.New(http.StatusText(http.StatusBadRequest))):
			return handler.HandleError(ctx, 400, err)
		default:
			return handler.HandleError(ctx, 500, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get an event by ID
// @Description Get details of a specific event by its ID
// @Tags events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} model.Response[model.EventResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /events/{id} [get]
func (h *EventHandlerImpl) GetEventByID(ctx echo.Context) error {
	request := new(model.GetEventRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.EventService.GetEventByID(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get event by id: %v", err)
		switch err {
		case domainErrors.ErrNotFound:
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get all events
// @Description Get a paginated list of all events
// @Tags events
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(name, date, time)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.EventResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /events [get]
func (h *EventHandlerImpl) GetAllEvents(ctx echo.Context) error {
	request := new(model.EventsRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.EventService.GetEvents(ctx.Request().Context(), request)
	if err != nil {
		switch err {
		case domainErrors.ErrValidation, domainErrors.ErrBadRequest:
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case domainErrors.ErrNotFound:
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// @Summary Search events
// @Description Search events with the provided query parameters
// @Tags events
// @Produce json
// @Param name query string false "Event name"
// @Param description query string false "Event description"
// @Param date query string false "Event date"
// @Param time query string false "Event time"
// @Param venue_id query int false "Venue ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(name, date, time)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.EventResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /events/search [get]
func (h *EventHandlerImpl) SearchEvents(ctx echo.Context) error {
	request := new(model.EventSearchRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.EventService.SearchEvents(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to search events: %v", err)
		switch {
		case errors.Is(err, errors.New(http.StatusText(http.StatusBadRequest))):
			return handler.HandleError(ctx, 400, err)
		default:
			return handler.HandleError(ctx, 500, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
