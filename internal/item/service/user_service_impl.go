package service

import (
	"context"
	"strconv"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type itemServiceImpl struct {
	repo repository.ItemRepository
}

// CreateCategory implements ItemService
func (s *itemServiceImpl) CreateCategory(body dto.CategoryRequest, ctx context.Context) (uint, error) {
	category := body.ToModel()
	err := s.repo.CreateCategory(category, ctx)
	if err != nil {
		return 0, err
	}
	return category.ID, err
}

// CreateItem implements ItemService
func (s *itemServiceImpl) CreateItem(body dto.ItemRequest, ctx context.Context) (uint, error) {
	item := body.ToModel()
	err := s.repo.CreateItem(item, ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, err
}

// DeleteCategory implements ItemService
func (s *itemServiceImpl) DeleteCategory(id string, ctx context.Context) error {
	panic("unimplemented")
}

// DeleteItem implements ItemService
func (s *itemServiceImpl) DeleteItem(id string, ctx context.Context) error {
	panic("unimplemented")
}

// FindCategories implements ItemService
func (s *itemServiceImpl) FindCategories(ctx context.Context) (dto.CategoriesResponse, error) {
	categories, err := s.repo.FindCategories(ctx)
	if err != nil {
		return nil, err
	}
	var categoriesResponse dto.CategoriesResponse
	categoriesResponse.FromModel(categories)
	return categoriesResponse, nil
}

// FindItem implements ItemService
func (s *itemServiceImpl) FindItem(id string, ctx context.Context) (*dto.ItemResponse, error) {
	panic("unimplemented")
}

// FindItemByCategory implements ItemService
func (s *itemServiceImpl) FindItemsByCategory(categoryId string, ctx context.Context) (dto.ItemsResponse, error) {
	id, err := strconv.Atoi(categoryId)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	items, err := s.repo.FindItemsByCategory(uint(id), ctx)
	if err != nil {
		return nil, err
	}
	var itemsResponse dto.ItemsResponse
	itemsResponse.FromModel(items)
	return itemsResponse, nil
}

// Findtems implements ItemService
func (s *itemServiceImpl) FindItems(ctx context.Context) (dto.ItemsResponse, error) {
	items, err := s.repo.FindItems(ctx)
	if err != nil {
		return nil, err
	}
	var itemsResponse dto.ItemsResponse
	itemsResponse.FromModel(items)
	return itemsResponse, nil
}

// UpdateCategory implements ItemService
func (s *itemServiceImpl) UpdateCategory(id string, body dto.CategoryRequest, ctx context.Context) error {
	panic("unimplemented")
}

// UpdateItem implements ItemService
func (s *itemServiceImpl) UpdateItem(id string, body dto.ItemRequest, ctx context.Context) error {
	panic("unimplemented")
}

func NewItemService(repository repository.ItemRepository) ItemService {
	return &itemServiceImpl{
		repo: repository,
	}
}
