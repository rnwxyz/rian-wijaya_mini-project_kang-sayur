package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	"github.com/stretchr/testify/mock"
)

type OrderServiceMock struct {
	mock.Mock
}

func (b *OrderServiceMock) CreateOrder(body dto.OrderRequest, userId string, ctx context.Context) (*dto.NewOrder, error) {
	args := b.Called()
	return args.Get(0).(*dto.NewOrder), args.Error(1)
}

func (b *OrderServiceMock) FindAllOrders(ctx context.Context) (dto.OrdersResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.OrdersResponse), args.Error(1)
}

func (b *OrderServiceMock) FindOrder(userId string, ctx context.Context) (dto.OrdersResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.OrdersResponse), args.Error(1)
}

func (b *OrderServiceMock) FindOrderDetail(userId string, orderId string, ctx context.Context) (*dto.OrderWithDetailResponse, error) {
	args := b.Called()
	return args.Get(0).(*dto.OrderWithDetailResponse), args.Error(1)
}

func (b *OrderServiceMock) CencelOder(orderId string, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderServiceMock) OrderReady(orderId string, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderServiceMock) TakeOrder(body dto.TakeOrder, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderServiceMock) SetOrderStatus(orderId uuid.UUID, status string, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}
