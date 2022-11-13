package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
	"github.com/stretchr/testify/mock"
)

type ItemServiceMock struct {
	mock.Mock
}

func (b *ItemServiceMock) CreateItem(body dto.ItemRequest, ctx context.Context) (uint, error) {
	args := b.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (b *ItemServiceMock) UpdateItem(id string, body dto.ItemRequest, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *ItemServiceMock) FindItems(ctx context.Context) (dto.ItemsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.ItemsResponse), args.Error(1)
}

func (b *ItemServiceMock) FindItemsByCategory(categoryId string, ctx context.Context) (dto.ItemsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.ItemsResponse), args.Error(1)
}

func (b *ItemServiceMock) CreateCategory(body dto.CategoryRequest, ctx context.Context) (uint, error) {
	args := b.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (b *ItemServiceMock) FindCategories(ctx context.Context) (dto.CategoriesResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.CategoriesResponse), args.Error(1)
}
