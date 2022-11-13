package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
	ism "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/service/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	mm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/middleware/mock"
	vm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/validator/mock"
	"github.com/stretchr/testify/suite"
)

type suiteItemController struct {
	suite.Suite
	itemServiceMock *ism.ItemServiceMock
	JWTServiceMock  *mm.MockJWTService
	validatorMock   *vm.CustomValidatorMock
	itemController  *itemController
	echoNew         *echo.Echo
}

func (s *suiteItemController) SetupSuit() {
	s.itemServiceMock = new(ism.ItemServiceMock)
	s.JWTServiceMock = new(mm.MockJWTService)
	s.validatorMock = new(vm.CustomValidatorMock)
	s.itemController = NewItemController(s.itemServiceMock, s.JWTServiceMock)
	s.echoNew = echo.New()
	s.echoNew.Validator = s.validatorMock
}

func (s *suiteItemController) TearDown() {
	s.itemServiceMock = nil
	s.JWTServiceMock = nil
	s.itemController = nil
	s.validatorMock = nil
	s.echoNew = nil
}

func (s *suiteItemController) TestCreateItem() {
	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		Body           map[string]interface{}
		JWTReturn      jwt.MapClaims
		ValidatorErr   error
		CreateItemErr  error
		CreateItemRes  uint
	}{
		{
			Name:           "success create",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"id":      float64(1),
				"message": "new item success created",
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  nil,
			CreateItemErr: nil,
			CreateItemRes: 1,
		},
		{
			Name:           "forbidden",
			ExpectedStatus: 403,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrPermission.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(100),
			},
			ValidatorErr:  nil,
			CreateItemErr: nil,
			CreateItemRes: 0,
		},
		{
			Name:           "invalid body type",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         "aa",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  nil,
			CreateItemErr: nil,
			CreateItemRes: 0,
		},
		{
			Name:           "bad body request",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  nil,
			CreateItemErr: customerrors.ErrBadRequestBody,
			CreateItemRes: 0,
		},
		{
			Name:           "item duplicate",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrDuplicateData.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  nil,
			CreateItemErr: customerrors.ErrDuplicateData,
			CreateItemRes: 0,
		},
		{
			Name:           "validator error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("validator error").Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  errors.New("validator error"),
			CreateItemErr: nil,
			CreateItemRes: 0,
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal errror").Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
				"category_id": 1,
				"qty":         1,
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:  nil,
			CreateItemErr: errors.New("internal errror"),
			CreateItemRes: 0,
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/items")

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.itemServiceMock.On("CreateItem").Return(v.CreateItemRes, v.CreateItemErr)
			s.validatorMock.On("Validate").Return(v.ValidatorErr)

			err = s.itemController.CreateItem(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}

func (s *suiteItemController) TestCreateCategory() {
	testCase := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedResult    map[string]interface{}
		Body              map[string]interface{}
		JWTReturn         jwt.MapClaims
		ValidatorErr      error
		CreateCategoryErr error
		CreateCategoryRes uint
	}{
		{
			Name:           "success create",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"id":      float64(1),
				"message": "new category success created",
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:      nil,
			CreateCategoryErr: nil,
			CreateCategoryRes: 1,
		},
		{
			Name:           "forbidden",
			ExpectedStatus: 403,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrPermission.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(100),
			},
			ValidatorErr:      nil,
			CreateCategoryErr: nil,
			CreateCategoryRes: 0,
		},
		{
			Name:           "bad body request",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:      nil,
			CreateCategoryErr: customerrors.ErrBadRequestBody,
			CreateCategoryRes: 0,
		},
		{
			Name:           "item duplicate",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrDuplicateData.Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:      nil,
			CreateCategoryErr: customerrors.ErrDuplicateData,
			CreateCategoryRes: 0,
		},
		{
			Name:           "validator error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("validator error").Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:      errors.New("validator error"),
			CreateCategoryErr: nil,
			CreateCategoryRes: 0,
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal errror").Error(),
			},
			Body: map[string]interface{}{
				"name":        "item",
				"description": "test",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:      nil,
			CreateCategoryErr: errors.New("internal errror"),
			CreateCategoryRes: 0,
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/items/categories")

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.itemServiceMock.On("CreateCategory").Return(v.CreateCategoryRes, v.CreateCategoryErr)
			s.validatorMock.On("Validate").Return(v.ValidatorErr)

			err = s.itemController.CreateCategory(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}
func (s *suiteItemController) TestGetItems() {
	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		FindItemsErr   error
		FindItemsRes   dto.ItemsResponse
	}{
		{
			Name:           "find items success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"category_name": "",
						"description":   "",
						"id":            float64(1),
						"name":          "item",
						"price":         float64(0),
						"qty":           float64(0),
					},
				},
				"message": "get items success",
			},
			FindItemsErr: nil,
			FindItemsRes: dto.ItemsResponse{
				{
					ID:   1,
					Name: "item",
				},
			},
		},
		{
			Name:           "error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error").Error(),
			},
			FindItemsErr: errors.New("error"),
			FindItemsRes: dto.ItemsResponse{},
		},
		{
			Name:           "nil items",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data":    []interface{}{},
				"message": "get items success",
			},
			FindItemsErr: nil,
			FindItemsRes: dto.ItemsResponse{},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/items")

			// define mock
			s.itemServiceMock.On("FindItems").Return(v.FindItemsRes, v.FindItemsErr)

			err := s.itemController.GetItems(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}
func (s *suiteItemController) TestGetCategories() {
	testCase := []struct {
		Name            string
		ExpectedStatus  int
		ExpectedResult  map[string]interface{}
		FindCategoryErr error
		FindCategoryRes dto.CategoriesResponse
	}{
		{
			Name:           "find categoris success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"name":        "category",
						"description": "test",
					},
				},
				"message": "get categories success",
			},
			FindCategoryErr: nil,
			FindCategoryRes: dto.CategoriesResponse{
				{
					ID:          1,
					Name:        "category",
					Description: "test",
				},
			},
		},
		{
			Name:           "error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error").Error(),
			},
			FindCategoryErr: errors.New("error"),
			FindCategoryRes: dto.CategoriesResponse{},
		},
		{
			Name:           "nil categories",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data":    []interface{}{},
				"message": "get categories success",
			},
			FindCategoryErr: nil,
			FindCategoryRes: dto.CategoriesResponse{},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/items")

			// define mock
			s.itemServiceMock.On("FindCategories").Return(v.FindCategoryRes, v.FindCategoryErr)

			err := s.itemController.GetCategories(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}
func (s *suiteItemController) TestGetItemByCategory() {
	testCase := []struct {
		Name                   string
		ExpectedStatus         int
		ParamId                string
		ExpectedResult         map[string]interface{}
		FindItemsByCategoryErr error
		FindItemsByCategoryRes dto.ItemsResponse
	}{
		{
			Name:           "find items success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"category_name": "",
						"description":   "",
						"id":            float64(1),
						"name":          "item",
						"price":         float64(0),
						"qty":           float64(0),
					},
				},
				"message": "get items success",
			},
			ParamId:                "1",
			FindItemsByCategoryErr: nil,
			FindItemsByCategoryRes: dto.ItemsResponse{
				{
					ID:   1,
					Name: "item",
				},
			},
		},
		{
			Name:           "bad request",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			ParamId:                "1",
			FindItemsByCategoryErr: customerrors.ErrBadRequestBody,
			FindItemsByCategoryRes: dto.ItemsResponse{},
		},
		{
			Name:           "internal error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error").Error(),
			},
			ParamId:                "1",
			FindItemsByCategoryErr: errors.New("error"),
			FindItemsByCategoryRes: dto.ItemsResponse{},
		},
		{
			Name:           "nil items",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data":    []interface{}{},
				"message": "get items success",
			},
			FindItemsByCategoryErr: nil,
			FindItemsByCategoryRes: dto.ItemsResponse{},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/items/category/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamId)

			// define mock
			s.itemServiceMock.On("FindItemsByCategory").Return(v.FindItemsByCategoryRes, v.FindItemsByCategoryErr)

			err := s.itemController.GetItemsByCategory(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}

func TestItemController(t *testing.T) {
	suite.Run(t, new(suiteItemController))
}
