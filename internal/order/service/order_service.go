package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
)

type OrderService interface {
	CreateOrder(body dto.OrderRequest, userId string, ctx context.Context) (uuid.UUID, error)
	FindAllOrders(ctx context.Context) (dto.OrdersResponse, error)
	FindOrder(userId string, ctx context.Context) (dto.OrdersResponse, error)
	FindOrderDetail(userId string, orderId string, ctx context.Context) (*dto.OrderWithDetailResponse, error)
	CencelOder(orderId string, ctx context.Context) error
	OrderReady(orderId string, ctx context.Context) error
	TakeOrder(code string, ctx context.Context) error
}
