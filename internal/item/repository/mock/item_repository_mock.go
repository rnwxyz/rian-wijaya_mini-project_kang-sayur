package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type ItemRepositoryMock struct {
	mock.Mock
}

func (b *ItemRepositoryMock) CreateItem(item *model.Item, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *ItemRepositoryMock) UpdateItem(item *model.Item, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *ItemRepositoryMock) FindItems(ctx context.Context) ([]model.Item, error) {
	args := b.Called()
	return args.Get(0).([]model.Item), args.Error(1)
}

func (b *ItemRepositoryMock) FindItemById(item *model.Item, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *ItemRepositoryMock) FindItemsByCategory(categoryId uint, ctx context.Context) ([]model.Item, error) {
	args := b.Called()
	return args.Get(0).([]model.Item), args.Error(1)
}

func (b *ItemRepositoryMock) CreateCategory(category *model.Category, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *ItemRepositoryMock) FindCategories(ctx context.Context) ([]model.Category, error) {
	args := b.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}
