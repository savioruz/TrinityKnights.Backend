package order

import (
	"context"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/order"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
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
	PaymentService  payment.PaymentService
}

func NewOrderServiceImpl(db *gorm.DB, cache *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, orderRepository order.OrderRepository, paymentService payment.PaymentService) *OrderServiceImpl {
	return &OrderServiceImpl{
		DB:              db,
		Cache:           cache,
		Log:             log,
		Validate:        validate,
		OrderRepository: orderRepository,
		PaymentService:  paymentService,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, request *model.OrderTicketRequest) (*model.OrderResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
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

	// Create order first
	order := &entity.Order{
		UserID:     request.UserID,
		Date:       time.Now(),
		TotalPrice: float64(request.Quantity) * 10.0, // @TODO: Change to event price
	}

	if err := s.OrderRepository.Create(tx, order); err != nil {
		s.Log.Errorf("failed to create order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create tickets
	for i := 0; i < request.Quantity; i++ {
		ticket := &entity.Ticket{
			EventID: request.EventID,
			OrderID: order.ID,
		}

		if err := tx.Create(ticket).Error; err != nil {
			s.Log.Errorf("failed to create ticket: %v", err)
			return nil, domainErrors.ErrInternalServer
		}
	}

	// Initialize payment
	paymentReq := &model.PaymentRequest{
		OrderID:       order.ID,
		Amount:        order.TotalPrice,
		PaymentMethod: request.PaymentMethod,
	}

	paymentResp, err := s.PaymentService.CreatePayment(ctx, paymentReq)
	if err != nil {
		s.Log.Errorf("failed to create payment: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Update order with payment ID
	order.PaymentID = paymentResp.ID
	if err := s.OrderRepository.Update(tx, order); err != nil {
		s.Log.Errorf("failed to update order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.OrderEntityToResponse(order), nil
}

func (s *OrderServiceImpl) UpdateOrder(ctx context.Context, request *model.UpdateOrderRequest) (*model.OrderResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var order entity.Order
	if err := s.OrderRepository.GetByID(tx, &order, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	order.PaymentID = request.PaymentID

	if err := s.OrderRepository.Update(tx, &order); err != nil {
		s.Log.Errorf("failed to update order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.OrderEntityToResponse(&order), nil
}

func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, request *model.GetOrderRequest) (*model.OrderResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	var order entity.Order
	if err := s.OrderRepository.GetByID(s.DB.WithContext(ctx), &order, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.OrderEntityToResponse(&order), nil
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
	if err := query.Offset(offset).Limit(request.Size).Find(&orders).Error; err != nil {
		s.Log.Errorf("failed to get orders: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(orders) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	return converter.OrdersToPaginatedResponse(orders, totalItems, request.Page, request.Size), nil
}
