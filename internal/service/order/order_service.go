package order

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, request *model.OrderTicketRequest) (*model.OrderResponse, error)
	UpdateOrder(ctx context.Context, request *model.UpdateOrderRequest) (*model.OrderResponse, error)
	GetOrderByID(ctx context.Context, request *model.GetOrderRequest) (*model.OrderResponse, error)
	GetOrders(ctx context.Context, request *model.OrdersRequest) (*model.Response[[]*model.OrderResponse], error)
}
