package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}

type orderController struct {
	service    service.OrderService
	jwtService JWTService
}

func NewOrderController(service service.OrderService, jwt JWTService) *orderController {
	return &orderController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *orderController) InitRoute(auth *echo.Group) {
	auth.POST("/order", u.CreateOrder)
	auth.GET("/order", u.GetOrder)
	auth.GET("/order/:id", u.GetOrderDetail)
	auth.PUT("/order/:id/cencel", u.CencelOrder)
	auth.PUT("/order/:id/ready", u.OrderReady)
}

func (u *orderController) CreateOrder(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)

	var orderBody dto.OrderRequest
	if err := c.Bind(&orderBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(orderBody); err != nil {
		return err
	}

	id, err := u.service.CreateOrder(orderBody.Order, userId, c.Request().Context())
	if err != nil {
		if err == utils.ErrInvalidId || err == utils.ErrQtyOrder || err == utils.ErrBadRequestBody {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "new create order success created",
		"id":      id,
	})
}

func (u *orderController) GetOrder(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)

	orders, err := u.service.FindOrder(userId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get orders success",
		"data":    orders,
	})
}

func (u *orderController) GetOrderDetail(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)
	orderId := c.Param("id")

	order, err := u.service.FindOrderDetail(userId, orderId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get orders success",
		"data":    order,
	})
}

func (u *orderController) CencelOrder(c echo.Context) error {
	orderId := c.Param("id")

	err := u.service.CencelOder(orderId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "order cencelled",
	})
}

func (u *orderController) OrderReady(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	orderId := c.Param("id")

	err := u.service.OrderReady(orderId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "order ready in checkpoint",
	})
}
