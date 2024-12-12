package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/event"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventServiceImpl struct {
	DB              *gorm.DB
	Cache           *cache.ImplCache
	Log             *logrus.Logger
	Validate        *validator.Validate
	EventRepository event.EventRepository
	helper          *helper.ContextHelper
}

func NewEventServiceImpl(db *gorm.DB, cache *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, eventRepository event.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{
		DB:              db,
		Cache:           cache,
		Log:             log,
		Validate:        validate,
		EventRepository: eventRepository,
		helper:          helper.NewContextHelper(),
	}
}

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04:05"
)

func (s *EventServiceImpl) CreateEvent(ctx context.Context, request *model.CreateEventRequest) (*model.EventResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	parsedDateTime, err := parseDateTime(request.Date, request.Time)
	if err != nil {
		s.Log.Errorf("failed to parse date time: %v", err)
		return nil, domainErrors.ErrValidation
	}

	data := &entity.Event{
		Name:        request.Name,
		Description: request.Description,
		Date:        parsedDateTime,
		Time:        helper.SQLTime(parsedDateTime),
		VenueID:     request.VenueID,
	}

	if err := s.EventRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create event: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.EventEntityToResponse(data), nil
}

func (s *EventServiceImpl) UpdateEvent(ctx context.Context, request *model.UpdateEventRequest) (*model.EventResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	parsedDateTime, err := parseDateTime(request.Date, request.Time)
	if err != nil {
		s.Log.Errorf("failed to parse date time: %v", err)
		return nil, domainErrors.ErrValidation
	}

	data := &entity.Event{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		Date:        parsedDateTime,
		Time:        helper.SQLTime(parsedDateTime),
		VenueID:     request.VenueID,
	}

	if err := s.EventRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update event: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.EventEntityToResponse(data), nil
}

func (s *EventServiceImpl) GetEventByID(ctx context.Context, request *model.GetEventRequest) (*model.EventResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	key := fmt.Sprintf("event:get:id:%d", request.ID)
	var data *model.EventResponse
	err := s.Cache.Get(key, &data)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		s.Log.Errorf("failed to get cache: %v", err)
	}

	if data == nil {
		tx := s.DB.WithContext(ctx)
		defer tx.Rollback()

		eventData := &entity.Event{}
		if err := s.EventRepository.GetByID(tx, eventData, request.ID); err != nil {
			return nil, domainErrors.ErrNotFound
		}

		response := converter.EventEntityToResponse(eventData)

		if err := s.Cache.Set(key, response, 5*time.Minute); err != nil {
			s.Log.Errorf("failed to set cache: %v", err)
		}

		return response, nil
	}

	return data, nil
}

func (s *EventServiceImpl) GetEvents(ctx context.Context, request *model.EventsRequest) (*model.Response[[]*model.EventResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.EventQueryOptions{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  request.Sort,
		Order: request.Order,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("event:get:page:%d:size:%d:sort:%s:order:%s", opts.Page, opts.Size, opts.Sort, opts.Order)
	var cacheResponse model.Response[[]*model.EventResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	db := s.DB.WithContext(ctx)

	var events []entity.Event
	totalItems, err := s.EventRepository.GetPaginated(db, &events, opts)
	if err != nil {
		s.Log.Errorf("failed to get events: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if totalItems == 0 || len(events) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.EventsToPaginatedResponse(events, totalItems, opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to cache response: %v", err)
	}

	return response, nil
}

func (s *EventServiceImpl) SearchEvents(ctx context.Context, request *model.EventSearchRequest) (*model.Response[[]*model.EventResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.EventQueryOptions{
		Page:        request.Page,
		Size:        request.Size,
		Name:        &request.Name,
		Description: &request.Description,
		Date:        &request.Date,
		Time:        &request.Time,
		VenueID:     &request.VenueID,
		Sort:        request.Sort,
		Order:       request.Order,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("event:search:%s:%s:%s:%s:%d:p%d:s%d:sort%s:%s",
		helper.StringOrEmpty(opts.Name),
		helper.StringOrEmpty(opts.Description),
		helper.StringOrEmpty(opts.Date),
		helper.StringOrEmpty(opts.Time),
		helper.UintOrZero(opts.VenueID),
		opts.Page,
		opts.Size,
		opts.Sort,
		opts.Order,
	)

	var cacheResponse model.Response[[]*model.EventResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	var events []entity.Event
	totalItems, err := s.EventRepository.GetPaginated(s.DB.WithContext(ctx), &events, opts)
	if err != nil {
		s.Log.Errorf("failed to search events: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if totalItems == 0 || len(events) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.EventsToPaginatedResponse(events, totalItems, opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to cache search results: %v", err)
	}

	return response, nil
}

func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	// Parse date
	date, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Parse time
	timeVal, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}

	// Combine date and time
	return time.Date(
		date.Year(), date.Month(), date.Day(),
		timeVal.Hour(), timeVal.Minute(), timeVal.Second(),
		0, time.Local,
	), nil
}
