package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}
type itemController struct {
	service    service.ItemService
	jwtService JWTService
}

func NewItemController(service service.ItemService, jwt JWTService) *itemController {
	return &itemController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *itemController) InitRoute(auth *echo.Group) {
	auth.POST("/item", u.CreateItem)
	auth.GET("/item", u.GetItems)
	auth.PUT("/item/:id", u.UpdateItem)
	auth.POST("/item/category", u.CreateCategory)
	auth.GET("/item/category", u.GetCategories)
	auth.GET("/item/category/:category_id", u.GetItemsByCategory)
}

func (u *itemController) CreateItem(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	var itemBody dto.ItemRequest
	if err := c.Bind(&itemBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(itemBody); err != nil {
		return err
	}
	id, err := u.service.CreateItem(itemBody, c.Request().Context())
	if err != nil {
		if err == utils.ErrBadRequestBody || err == utils.ErrDuplicateData {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "new item success created",
		"id":      id,
	})
}

func (u *itemController) CreateCategory(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	var categoryBody dto.CategoryRequest
	if err := c.Bind(&categoryBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(categoryBody); err != nil {
		return err
	}
	id, err := u.service.CreateCategory(categoryBody, c.Request().Context())
	if err != nil {
		if err == utils.ErrBadRequestBody || err == utils.ErrDuplicateData {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "new category success created",
		"id":      id,
	})
}

func (u *itemController) UpdateItem(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	id := c.Param("id")
	var itemBody dto.ItemRequest
	if err := c.Bind(&itemBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(itemBody); err != nil {
		return err
	}
	err := u.service.UpdateItem(id, itemBody, c.Request().Context())
	if err != nil {
		if err == utils.ErrBadRequestBody || err == utils.ErrDuplicateData {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update item",
	})
}

func (u *itemController) GetItems(c echo.Context) error {
	items, err := u.service.FindItems(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get items success",
		"data":    items,
	})
}

func (u *itemController) GetCategories(c echo.Context) error {
	categories, err := u.service.FindCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get categories success",
		"data":    categories,
	})
}

func (u *itemController) GetItemsByCategory(c echo.Context) error {
	paramId := c.Param("category_id")
	items, err := u.service.FindItemsByCategory(paramId, c.Request().Context())
	if err != nil {
		if err == utils.ErrBadRequestBody {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get items success",
		"data":    items,
	})
}