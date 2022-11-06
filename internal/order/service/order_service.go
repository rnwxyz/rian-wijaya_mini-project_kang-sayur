package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
)

type OrderService interface {
	CreateOrder(body dto.OrderDetailsRequest, userId string, ctx context.Context) (uuid.UUID, error)
	FindOrder(userId string, ctx context.Context) (dto.OrdersResponse, error)
	FindOrderDetail(userId string, orderId string, ctx context.Context) (*dto.OrderWithDetailResponse, error)
	CencelOder(orderId string, ctx context.Context) error
}
