package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type OrderRepository interface {
	CreateOrder(order *model.Order, ctx context.Context) error
	FindAllOrders(ctx context.Context) ([]model.Order, error)
	FindOrder(userId uuid.UUID, ctx context.Context) ([]model.Order, error)
	FindOrderDetail(order *model.Order, ctx context.Context) error
	CencelOrder(orderId uuid.UUID, ctx context.Context) error
	OrderReady(orderId uuid.UUID, ctx context.Context) error
	OrderDone(orderId uuid.UUID, ctx context.Context) error
	OrderWaiting(orderId uuid.UUID, ctx context.Context) error
	InitStatusOrder() error
}
