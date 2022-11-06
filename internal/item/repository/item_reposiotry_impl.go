package repository

import (
	"context"
	"strings"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

type itemRepositoryImpl struct {
	db *gorm.DB
}

// CreateCategory implements ItemRepository
func (r *itemRepositoryImpl) CreateCategory(category *model.Category, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrDuplicateData
		}
		return err
	}
	return nil
}

// CreateItem implements ItemRepository
func (r *itemRepositoryImpl) CreateItem(item *model.Item, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(item).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrDuplicateData
		}
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return utils.ErrBadRequestBody
		}
		return err
	}
	return nil
}

// DeleteItem implements ItemRepository
func (r *itemRepositoryImpl) DeleteItem(item *model.Item, ctx context.Context) error {
	panic("unimplemented")
}

// FindCategories implements ItemRepository
func (r *itemRepositoryImpl) FindCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, err
}

// FindItemById implements ItemRepository
func (r *itemRepositoryImpl) FindItemById(item *model.Item, ctx context.Context) error {
	err := r.db.WithContext(ctx).First(item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

// FindItems implements ItemRepository
func (r *itemRepositoryImpl) FindItems(ctx context.Context) ([]model.Item, error) {
	var items []model.Item
	err := r.db.WithContext(ctx).Preload("Category").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, err
}

// FindItemByCategory implements ItemRepository
func (r *itemRepositoryImpl) FindItemsByCategory(categoryId uint, ctx context.Context) ([]model.Item, error) {
	var items []model.Item
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryId).Preload("Category").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, err
}

// UpdateItem implements ItemRepository
func (r *itemRepositoryImpl) UpdateItem(item *model.Item, ctx context.Context) error {
	res := r.db.WithContext(ctx).Model(&model.Item{}).Where("id = ?", item.ID).Updates(&model.Item{
		Name:        item.Name,
		Description: item.Description,
		Qty:         item.Qty,
		Price:       item.Price,
		CategoryID:  item.CategoryID,
	})

	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "Duplicate entry") {
			return utils.ErrDuplicateData
		}
		if strings.Contains(res.Error.Error(), "Cannot add or update a child row") {
			return utils.ErrBadRequestBody
		}
		return res.Error
	}
	if res.RowsAffected == 0 {
		return utils.ErrInvalidId
	}
	return nil
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepositoryImpl{
		db: db,
	}
}
