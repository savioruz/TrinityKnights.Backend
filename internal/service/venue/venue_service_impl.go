package venue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/venue"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VenueServiceImpl struct {
	DB              *gorm.DB
	Cache           *cache.ImplCache
	Log             *logrus.Logger
	Validate        *validator.Validate
	VenueRepository venue.VenueRepository
	helper          *helper.ContextHelper
}

func NewVenueServiceImpl(db *gorm.DB, cache *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, venueRepository venue.VenueRepository) *VenueServiceImpl {
	return &VenueServiceImpl{
		DB:              db,
		Cache:           cache,
		Log:             log,
		Validate:        validate,
		VenueRepository: venueRepository,
		helper:          helper.NewContextHelper(),
	}
}

func (s *VenueServiceImpl) CreateVenue(ctx context.Context, request *model.CreateVenueRequest) (*model.VenueResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	data := &entity.Venue{
		Name:     request.Name,
		Address:  request.Address,
		Capacity: request.Capacity,
		City:     request.City,
		State:    request.State,
		Zip:      request.Zip,
	}

	if err := s.VenueRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create venue: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.VenueEntityToResponse(data), nil
}

func (s *VenueServiceImpl) UpdateVenue(ctx context.Context, request *model.UpdateVenueRequest) (*model.VenueResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	data := &entity.Venue{
		ID:       request.ID,
		Name:     request.Name,
		Address:  request.Address,
		Capacity: request.Capacity,
		City:     request.City,
		State:    request.State,
		Zip:      request.Zip,
	}

	if err := s.VenueRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update venue: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.VenueEntityToResponse(data), nil
}

func (s *VenueServiceImpl) GetVenueByID(ctx context.Context, request *model.GetVenueRequest) (*model.VenueResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	key := fmt.Sprintf("venue:get:id:%d", request.ID)
	var data *model.VenueResponse
	err := s.Cache.Get(key, &data)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		s.Log.Errorf("failed to get cache: %v", err)
	}

	if data == nil {
		tx := s.DB.WithContext(ctx)
		defer tx.Rollback()

		venueData := &entity.Venue{}
		if err := s.VenueRepository.GetByID(tx, venueData, request.ID); err != nil {
			return nil, domainErrors.ErrNotFound
		}

		response := converter.VenueEntityToResponse(venueData)

		if err := s.Cache.Set(key, response, 5*time.Minute); err != nil {
			s.Log.Errorf("failed to set cache: %v", err)
		}

		return response, nil
	}

	return data, nil
}

func (s *VenueServiceImpl) GetVenues(ctx context.Context, request *model.VenuesRequest) (*model.Response[[]*model.VenueResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.VenueQueryOptions{
		Page: request.Page,
		Size: request.Size,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("venue:get:page:%d:size:%d", opts.Page, opts.Size)
	var cacheResponse model.Response[[]*model.VenueResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	db := s.DB.WithContext(ctx)

	var venues []entity.Venue
	totalItems, err := s.VenueRepository.GetPaginated(db, &venues, opts)
	if err != nil {
		s.Log.Errorf("failed to get venues: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if totalItems == 0 || len(venues) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.VenuesToPaginatedResponse(venues, totalItems, opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to cache response: %v", err)
	}

	return response, nil
}

func (s *VenueServiceImpl) SearchVenues(ctx context.Context, request *model.VenueSearchRequest) (*model.Response[[]*model.VenueResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.VenueQueryOptions{
		Page:     request.Page,
		Size:     request.Size,
		Name:     &request.Name,
		Address:  &request.Address,
		Capacity: &request.Capacity,
		City:     &request.City,
		State:    &request.State,
		Zip:      &request.Zip,
		Sort:     request.Sort,
		Order:    request.Order,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("venue:search:%s:%s:%d:%s:%s:%s:p%d:s%d:sort%s:%s",
		helper.StringOrEmpty(opts.Name),
		helper.StringOrEmpty(opts.Address),
		helper.IntOrZero(opts.Capacity),
		helper.StringOrEmpty(opts.City),
		helper.StringOrEmpty(opts.State),
		helper.StringOrEmpty(opts.Zip),
		opts.Page,
		opts.Size,
		opts.Sort,
		opts.Order,
	)

	var cacheResponse model.Response[[]*model.VenueResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	var venues []entity.Venue
	totalItems, err := s.VenueRepository.GetPaginated(s.DB.WithContext(ctx), &venues, opts)
	if err != nil {
		s.Log.Errorf("failed to search venues: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if totalItems == 0 || len(venues) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.VenuesToPaginatedResponse(venues, totalItems, opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to cache search results: %v", err)
	}

	return response, nil
}
