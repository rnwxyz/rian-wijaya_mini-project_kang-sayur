package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	itemRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
)

type suiteItemService struct {
	suite.Suite
	itemRepositoryMock *itemRepositoryMock.ItemRepositoryMock
	itemService        ItemService
}

func newItemService(repository repository.ItemRepository) ItemService {
	return &itemServiceImpl{
		repo: repository,
	}
}

func (s *suiteItemService) SetupSuit() {
	s.itemRepositoryMock = new(itemRepositoryMock.ItemRepositoryMock)
	s.itemService = newItemService(s.itemRepositoryMock)
}

func (s *suiteItemService) TearDown() {
	s.itemRepositoryMock = nil
	s.itemService = nil
}

func (s *suiteItemService) TestCreateCategory() {
	testCase := []struct {
		Name              string
		Body              dto.CategoryRequest
		ExpectedErr       error
		CreateCategoryErr error
	}{
		{
			Name: "success",
			Body: dto.CategoryRequest{
				Name:        "category",
				Description: "test",
			},
			ExpectedErr:       nil,
			CreateCategoryErr: nil,
		},
		{
			Name: "error",
			Body: dto.CategoryRequest{
				Name:        "category",
				Description: "test",
			},
			ExpectedErr:       errors.New("error"),
			CreateCategoryErr: errors.New("error"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.itemRepositoryMock.On("CreateCategory").Return(v.CreateCategoryErr)

			_, err := s.itemService.CreateCategory(v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteItemService) TestCreateItem() {
	testCase := []struct {
		Name          string
		Body          dto.ItemRequest
		ExpectedErr   error
		CreateItemErr error
	}{
		{
			Name: "success",
			Body: dto.ItemRequest{
				Name:        "category",
				Description: "test",
			},
			ExpectedErr:   nil,
			CreateItemErr: nil,
		},
		{
			Name: "error",
			Body: dto.ItemRequest{
				Name:        "category",
				Description: "test",
			},
			ExpectedErr:   errors.New("error"),
			CreateItemErr: errors.New("error"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.itemRepositoryMock.On("CreateItem").Return(v.CreateItemErr)

			_, err := s.itemService.CreateItem(v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteItemService) TestFindCategory() {
	testCase := []struct {
		Name            string
		ExpectedErr     error
		ExpectedRes     dto.CategoriesResponse
		FindCategoryErr error
		FindCategoryRes []model.Category
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: dto.CategoriesResponse{
				{ID: 1, Name: "category"},
			},
			FindCategoryErr: nil,
			FindCategoryRes: []model.Category{
				{
					ID:   1,
					Name: "category",
				},
			},
		},
		{
			Name:            "error",
			ExpectedErr:     errors.New("error"),
			ExpectedRes:     dto.CategoriesResponse(nil),
			FindCategoryErr: errors.New("error"),
			FindCategoryRes: []model.Category{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.itemRepositoryMock.On("FindCategories").Return(v.FindCategoryRes, v.FindCategoryErr)

			res, err := s.itemService.FindCategories(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteItemService) TestFindItemByCategory() {
	testCase := []struct {
		Name              string
		ExpectedErr       error
		ExpectedRes       dto.ItemsResponse
		CategoryId        string
		FindByCategoryErr error
		FindByCategoryRes []model.Item
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: dto.ItemsResponse{
				{ID: 1, Name: "item"},
			},
			CategoryId:        "1",
			FindByCategoryErr: nil,
			FindByCategoryRes: []model.Item{
				{
					ID:         1,
					Name:       "item",
					CategoryID: 1,
				},
			},
		},
		{
			Name:              "error",
			ExpectedErr:       errors.New("error"),
			ExpectedRes:       dto.ItemsResponse(nil),
			CategoryId:        "1",
			FindByCategoryErr: errors.New("error"),
			FindByCategoryRes: []model.Item{},
		},
		{
			Name:              "invalid id",
			ExpectedErr:       customerrors.ErrInvalidId,
			ExpectedRes:       dto.ItemsResponse(nil),
			CategoryId:        "ewe",
			FindByCategoryErr: nil,
			FindByCategoryRes: []model.Item{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.itemRepositoryMock.On("FindItemsByCategory").Return(v.FindByCategoryRes, v.FindByCategoryErr)

			res, err := s.itemService.FindItemsByCategory(v.CategoryId, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteItemService) TestFindItems() {
	testCase := []struct {
		Name         string
		ExpectedErr  error
		ExpectedRes  dto.ItemsResponse
		FindItemsErr error
		FindItemsRes []model.Item
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: dto.ItemsResponse{
				{ID: 1, Name: "item"},
			},
			FindItemsErr: nil,
			FindItemsRes: []model.Item{
				{
					ID:         1,
					Name:       "item",
					CategoryID: 1,
				},
			},
		},
		{
			Name:         "error",
			ExpectedErr:  errors.New("error"),
			ExpectedRes:  dto.ItemsResponse(nil),
			FindItemsErr: errors.New("error"),
			FindItemsRes: []model.Item{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.itemRepositoryMock.On("FindItems").Return(v.FindItemsRes, v.FindItemsErr)

			res, err := s.itemService.FindItems(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteItemService(t *testing.T) {
	suite.Run(t, new(suiteItemService))
}
