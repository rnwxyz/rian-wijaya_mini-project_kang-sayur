package service

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
)

type ItemService interface {
	CreateItem(body dto.ItemRequest, ctx context.Context) (uint, error)
	UpdateItem(id string, body dto.ItemRequest, ctx context.Context) error
	FindItems(ctx context.Context) (dto.ItemsResponse, error)
	FindItemsByCategory(categoryId string, ctx context.Context) (dto.ItemsResponse, error)
	CreateCategory(body dto.CategoryRequest, ctx context.Context) (uint, error)
	FindCategories(ctx context.Context) (dto.CategoriesResponse, error)
}
