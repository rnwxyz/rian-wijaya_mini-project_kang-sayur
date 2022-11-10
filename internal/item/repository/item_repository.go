package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type ItemRepository interface {
	CreateItem(item *model.Item, ctx context.Context) error
	UpdateItem(item *model.Item, ctx context.Context) error
	FindItems(ctx context.Context) ([]model.Item, error)
	FindItemById(item *model.Item, ctx context.Context) error
	FindItemsByCategory(categoryId uint, ctx context.Context) ([]model.Item, error)
	CreateCategory(category *model.Category, ctx context.Context) error
	FindCategories(ctx context.Context) ([]model.Category, error)
}
