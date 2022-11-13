package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/service"
	tsm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/service/mock"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	mm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/middleware/mock"
	"github.com/stretchr/testify/suite"
)

type suiteTransactionController struct {
	suite.Suite
	transactionServiceMock *tsm.TransactionServiceMock
	JWTServiceMock         *mm.MockJWTService
	transactionController  *transactionController
}

func newTransactionControllerMock(service service.TransactionService, jwt JWTService) *transactionController {
	return &transactionController{
		service:    service,
		jwtService: jwt,
	}
}

func (s *suiteTransactionController) SetupTest() {
	s.transactionServiceMock = new(tsm.TransactionServiceMock)
	s.JWTServiceMock = new(mm.MockJWTService)
	s.transactionController = newTransactionControllerMock(s.transactionServiceMock, s.JWTServiceMock)
}

func (s *suiteTransactionController) TestTransactionNotification() {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name                 string
		Body                 map[string]interface{}
		ExpectedStatus       int
		CreateTransactionErr error
		ExpectedResult       map[string]interface{}
	}{
		{
			Name: "success",
			Body: map[string]interface{}{
				"transaction_id":     transactionId.String(),
				"order_id":           orderId.String(),
				"transaction_status": "pending",
				"transaction_time":   "2022-01-01 00:00:00",
				"signature_key":      "003592974b7cb5956dbfe36c000b",
				"payment_type":       "qris",
				"gross_amount":       "10000",
				"settlement_time":    "",
			},
			ExpectedStatus:       200,
			CreateTransactionErr: nil,
			ExpectedResult: map[string]interface{}{
				"message": "transaction success created or updated",
			},
		},
		{
			Name: "body type invalid",
			Body: map[string]interface{}{
				"transaction_id":     transactionId.String(),
				"order_id":           orderId.String(),
				"transaction_status": "pending",
				"transaction_time":   "2022-01-01 00:00:00",
				"signature_key":      "003592974b7cb5956dbfe36c000b",
				"payment_type":       "qris",
				"gross_amount":       10000,
				"settlement_time":    "",
			},
			ExpectedStatus:       400,
			CreateTransactionErr: nil,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
		},
		{
			Name:                 "empty body",
			Body:                 map[string]interface{}{},
			ExpectedStatus:       500,
			CreateTransactionErr: errors.New("internal server error"),
			ExpectedResult: map[string]interface{}{
				"message": "internal server error",
			},
		},
		{
			Name: "internal server error",
			Body: map[string]interface{}{
				"transaction_id":     transactionId.String(),
				"order_id":           orderId.String(),
				"transaction_status": "pending",
				"transaction_time":   "2022-01-01 00:00:00",
				"signature_key":      "003592974b7cb5956dbfe36c000b",
				"payment_type":       "qris",
				"gross_amount":       "10000",
				"settlement_time":    "",
			},
			ExpectedStatus:       500,
			CreateTransactionErr: errors.New("internal server error"),
			ExpectedResult: map[string]interface{}{
				"message": "internal server error",
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/transaction/notification")

			// define mock
			mock1 := s.transactionServiceMock.On("CreateTransaction").Return(v.CreateTransactionErr)

			err = s.transactionController.TransactionNotification(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			mock1.Unset()
		})
	}
}
func (s *suiteTransactionController) TestGetTransaction() {
	userId := uuid.New()
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedStatus     int
		FindTransactionErr error
		FindTransactionRes dto.TransactionsResponse
		JwtRes             jwt.MapClaims
		ExpectedResult     map[string]interface{}
	}{
		{
			Name:               "success",
			ExpectedStatus:     200,
			FindTransactionErr: nil,
			FindTransactionRes: dto.TransactionsResponse{
				{
					TransactionID:     transactionId.String(),
					OrderID:           orderId.String(),
					TransactionStatus: "pending",
					TransactionTime:   "2022-01-01 00:00:00",
					SignatureKey:      "003592974b7cb5956dbfe36c000b",
					PaymentType:       "qris",
					GrossAmount:       "10000",
					SettlementTime:    "",
				},
			},
			JwtRes: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ExpectedResult: map[string]interface{}{
				"message": "get transactions success",
				"data": []interface{}{
					map[string]interface{}{
						"transaction_id":     transactionId.String(),
						"order_id":           orderId.String(),
						"transaction_status": "pending",
						"transaction_time":   "2022-01-01 00:00:00",
						"signature_key":      "003592974b7cb5956dbfe36c000b",
						"payment_type":       "qris",
						"gross_amount":       "10000",
						"settlement_time":    "",
					},
				},
			},
		},
		{
			Name:               "internal server error",
			ExpectedStatus:     500,
			FindTransactionErr: errors.New("internal server error"),
			FindTransactionRes: nil,
			JwtRes: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ExpectedResult: map[string]interface{}{
				"message": "internal server error",
			},
		},
		{
			Name:               "nil transactions",
			ExpectedStatus:     200,
			FindTransactionErr: nil,
			FindTransactionRes: nil,
			JwtRes: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ExpectedResult: map[string]interface{}{
				"message": "get transactions success",
				"data":    nil,
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/transaction/notification")

			// define mock
			mock1 := s.transactionServiceMock.On("FindTransaction").Return(v.FindTransactionRes, v.FindTransactionErr)
			mock2 := s.JWTServiceMock.On("GetClaims").Return(v.JwtRes)

			err := s.transactionController.GetTransactions(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			mock1.Unset()
			mock2.Unset()
		})
	}
}

func (s *suiteTransactionController) TestInitRoute() {
	group := echo.New().Group("/api/v1")
	s.NotPanics(func() {
		s.transactionController.InitRoute(group, group)
	})
}

func TestTransactionController(t *testing.T) {
	suite.Run(t, new(suiteTransactionController))
}
