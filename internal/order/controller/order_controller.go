package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
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
	orders := auth.Group("/orders")
	orders.POST("", u.CreateOrder)
	orders.GET("", u.GetOrder)
	orders.GET("/:id", u.GetOrderDetail)
	orders.POST("/takeorder", u.TakeOrder)
	orders.PUT("/cencel/:id", u.CencelOrder)
	orders.PUT("/ready/:id", u.OrderReady)
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

	newOrder, err := u.service.CreateOrder(orderBody, userId, c.Request().Context())
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
		"data":    newOrder,
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

func (u *orderController) TakeOrder(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role != constants.Role_admin {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	var takeOrder dto.TakeOrder
	if err := c.Bind(&takeOrder); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrOrderCode.Error()})
	}
	if err := c.Validate(takeOrder); err != nil {
		return err
	}
	err := u.service.TakeOrder(takeOrder.Code, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrOrderCode.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success take order",
	})
}

func (u *orderController) OrderReady(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role != constants.Role_admin {
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
