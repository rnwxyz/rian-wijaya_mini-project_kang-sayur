package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
)

type OrderService interface {
	CreateOrder(body dto.OrderRequest, userId string, ctx context.Context) (*dto.NewOrder, error)
	FindAllOrders(ctx context.Context) (dto.OrdersResponse, error)
	FindOrder(userId string, ctx context.Context) (dto.OrdersResponse, error)
	FindOrderDetail(userId string, orderId string, ctx context.Context) (*dto.OrderWithDetailResponse, error)
	CencelOder(orderId string, ctx context.Context) error
	OrderReady(orderId string, ctx context.Context) error
	TakeOrder(body dto.TakeOrder, ctx context.Context) error
	SetOrderStatus(orderId uuid.UUID, status string, ctx context.Context) error
}
