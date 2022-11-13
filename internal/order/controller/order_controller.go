package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}

type QRCode interface {
	GenerateQRCode(hashCode string) ([]byte, error)
}

type orderController struct {
	service    service.OrderService
	jwtService JWTService
	qrCode     QRCode
}

func NewOrderController(service service.OrderService, jwt JWTService, qr QRCode) *orderController {
	return &orderController{
		service:    service,
		jwtService: jwt,
		qrCode:     qr,
	}
}

func (u *orderController) InitRoute(auth *echo.Group) {
	orders := auth.Group("/orders")
	orders.POST("", u.CreateOrder)
	orders.GET("", u.GetOrder)
	orders.GET("/:id", u.GetOrderDetail)
	orders.GET("/qr/:hash_code", u.GetQRCode)
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
			"message": customerrors.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(orderBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error()})
	}

	newOrder, err := u.service.CreateOrder(orderBody, userId, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrInvalidId || err == customerrors.ErrQtyOrder || err == customerrors.ErrBadRequestBody {
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
			"message": customerrors.ErrPermission.Error(),
		})
	}
	var takeOrder dto.TakeOrder
	if err := c.Bind(&takeOrder); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrOrderCode.Error()})
	}
	if err := c.Validate(takeOrder); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error()})
	}
	err := u.service.TakeOrder(takeOrder, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
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
			"message": customerrors.ErrPermission.Error(),
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

func (u *orderController) GetQRCode(c echo.Context) error {
	hashCode := c.Param("hash_code")
	qr, err := u.qrCode.GenerateQRCode(hashCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrGenerateQR.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get qrcode success",
		"data": map[string][]byte{
			"base64": qr,
		},
	})
}
