package venue

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/venue"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type VenueHandlerImpl struct {
	Log          *logrus.Logger
	VenueService venue.VenueService
}

func NewVenueHandler(log *logrus.Logger, venueService venue.VenueService) VenueHandler {
	return &VenueHandlerImpl{
		Log:          log,
		VenueService: venueService,
	}
}

// CreateVenue function is a handler to create a new venue
// @Summary Create a new venue @admin
// @Description Create a new venue with the provided details
// @Tags venues
// @Accept json
// @Produce json
// @Param request body model.CreateVenueRequest true "Venue details"
// @Success 201 {object} model.Response[model.VenueResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /venues [post]
func (h *VenueHandlerImpl) CreateVenue(ctx echo.Context) error {
	request := new(model.CreateVenueRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.VenueService.CreateVenue(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create venue: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// UpdateVenue function is a handler to update an existing venue
// @Summary Update an existing venue @admin
// @Description Update an existing venue with the provided details
// @Tags venues
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param request body model.UpdateVenueRequest true "Updated venue details"
// @Success 200 {object} model.Response[model.VenueResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /venues/{id} [put]
func (h *VenueHandlerImpl) UpdateVenue(ctx echo.Context) error {
	request := new(model.UpdateVenueRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.VenueService.UpdateVenue(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to update venue: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// GetVenueByID function is a handler to get a venue by its ID
// @Summary Get a venue by ID @admin
// @Description Get details of a specific venue by its ID
// @Tags venues
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} model.Response[model.VenueResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /venues/{id} [get]
func (h *VenueHandlerImpl) GetVenueByID(ctx echo.Context) error {
	request := new(model.GetVenueRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.VenueService.GetVenueByID(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get venue by id: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// GetAllVenues function is a handler to get all venues
// @Summary Get all venues @admin
// @Description Get a paginated list of all venues
// @Tags venues
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(name, capacity)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.VenueResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /venues [get]
func (h *VenueHandlerImpl) GetAllVenues(ctx echo.Context) error {
	request := new(model.VenuesRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.VenueService.GetVenues(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get venues: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// SearchVenues function is a handler to search venues
// @Summary Search venues @admin
// @Description Search venues with the provided query parameters
// @Tags venues
// @Produce json
// @Param id query int false "Venue ID"
// @Param name query string false "Venue name"
// @Param address query string false "Venue address"
// @Param capacity query int false "Venue capacity"
// @Param city query string false "Venue city"
// @Param state query string false "Venue state"
// @Param zip query string false "Venue zip"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(name, capacity)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.VenueResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /venues/search [get]
func (h *VenueHandlerImpl) SearchVenues(ctx echo.Context) error {
	request := new(model.VenueSearchRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.VenueService.SearchVenues(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to search venues: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrValidation):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
