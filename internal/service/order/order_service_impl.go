package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/order"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	DB              *gorm.DB
	Cache           *cache.ImplCache
	Log             *logrus.Logger
	Validate        *validator.Validate
	OrderRepository order.OrderRepository
	helper          *helper.ContextHelper
}

func NewOrderServiceImpl(db *gorm.DB, cacheImpl *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, orderRepository order.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{
		DB:              db,
		Cache:           cacheImpl,
		Log:             log,
		Validate:        validate,
		OrderRepository: orderRepository,
		helper:          helper.NewContextHelper(),
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, request *model.OrderTicketRequest) (*model.OrderResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	// Get user ID from JWT claims
	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		return nil, domainErrors.ErrUnauthorized
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Check if event exists
	var event entity.Event
	if err := tx.First(&event, request.EventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get event: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create data order
	data := &entity.Order{
		UserID:     claims.UserID,
		Date:       time.Now(),
		TotalPrice: float64(request.Quantity) * 10.0, // @TODO: Change to ticket price
	}

	if err := s.OrderRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create data: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create tickets
	for i := 0; i < request.Quantity; i++ {
		ticket := &entity.Ticket{
			EventID:    request.EventID,
			Price:      10.0, // @TODO: Change to ticket price
			Type:       "regular",
			SeatNumber: request.SeatNumber,
		}

		if err := tx.Create(ticket).Error; err != nil {
			s.Log.Errorf("failed to create ticket: %v", err)
			return nil, domainErrors.ErrInternalServer
		}
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.OrderEntityToResponse(data), nil
}

func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, request *model.GetOrderRequest) (*model.OrderResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	// Try to get from cache first
	key := fmt.Sprintf("order:get:id:%d", request.ID)
	var data *model.OrderResponse
	err := s.Cache.Get(key, &data)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		s.Log.Errorf("failed to get cache: %v", err)
	}

	if data == nil {
		tx := s.DB.WithContext(ctx)

		var dataOrder entity.Order
		if err := s.OrderRepository.GetByIDWithDetails(tx, &dataOrder, request.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, domainErrors.ErrNotFound
			}
			s.Log.Errorf("failed to get order: %v", err)
			return nil, domainErrors.ErrInternalServer
		}

		response := converter.OrderEntityToResponse(&dataOrder)

		// Cache the response
		if err := s.Cache.Set(key, response, 5*time.Minute); err != nil {
			s.Log.Errorf("failed to set cache: %v", err)
		}

		return response, nil
	}

	return data, nil
}

func (s *OrderServiceImpl) GetOrders(ctx context.Context, request *model.OrdersRequest) (*model.Response[[]*model.OrderResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	if request.Size <= 0 {
		request.Size = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("order:get:page:%d:size:%d:sort:%s:order:%s",
		request.Page, request.Size, request.Sort, request.Order)
	var cacheResponse model.Response[[]*model.OrderResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	var orders []entity.Order
	var totalItems int64

	query := s.DB.WithContext(ctx).Model(&entity.Order{})

	// Add sorting
	if request.Sort != "" && request.Order != "" {
		validSortFields := map[string]bool{
			"date":        true,
			"total_price": true,
			"created_at":  true,
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		if validSortFields[request.Sort] && validOrders[request.Order] {
			query = query.Order(fmt.Sprintf("%s %s", request.Sort, request.Order))
		} else {
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		s.Log.Errorf("failed to count orders: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Get paginated results
	offset := (request.Page - 1) * request.Size
	if err := query.Preload("Tickets").Offset(offset).Limit(request.Size).Find(&orders).Error; err != nil {
		s.Log.Errorf("failed to get orders: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(orders) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.OrdersToPaginatedResponse(orders, totalItems, request.Page, request.Size)

	// Cache the response
	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to cache response: %v", err)
	}

	return response, nil
}
