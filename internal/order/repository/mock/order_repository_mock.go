package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (b *OrderRepositoryMock) CreateOrder(order *model.Order, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	args := b.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func (b *OrderRepositoryMock) FindOrder(userId uuid.UUID, ctx context.Context) ([]model.Order, error) {
	args := b.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func (b *OrderRepositoryMock) FindOrderDetail(order *model.Order, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) CencelOrder(orderId uuid.UUID, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) FindOrderById(order *model.Order, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) OrderReady(orderId uuid.UUID, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) OrderDone(orderId uuid.UUID, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) OrderWaiting(orderId uuid.UUID, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *OrderRepositoryMock) InitStatusOrder() error {
	args := b.Called()
	return args.Error(0)
}
