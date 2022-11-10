package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/service"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}

type transactionController struct {
	service    service.TransactionService
	jwtService JWTService
}

func NewTransactionController(service service.TransactionService, jwt JWTService) *transactionController {
	return &transactionController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *transactionController) InitRoute(api, auth *echo.Group) {
	transactions := api.Group("/transactions")
	transactions.POST("/notification", u.TransactionNotification)

	transactionsWithAuth := auth.Group("/transactions")
	transactionsWithAuth.GET("", u.GetTransactions)
}

func (u *transactionController) TransactionNotification(c echo.Context) error {
	var transactionBody dto.TransactionRequest

	if err := c.Bind(&transactionBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error()})
	}

	err := u.service.CreateTransaction(transactionBody, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "transaction success created or updated",
	})
}

func (u *transactionController) GetTransactions(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)

	transactions, err := u.service.FindTransaction(userId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get transactions success",
		"data":    transactions,
	})
}
