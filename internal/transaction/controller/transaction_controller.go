package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/service"
)

type transactionController struct {
	service service.TransactionService
}

func NewTransactionController(service service.TransactionService) *transactionController {
	return &transactionController{
		service: service,
	}
}

func (u *transactionController) InitRoute(api *echo.Group) {
	api.POST("/transaction/notification", u.TransactionNotification)
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
