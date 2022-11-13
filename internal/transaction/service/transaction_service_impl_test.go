package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	_orderRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository/mock"
	os "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	_orderServiceMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	tr "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/repository"
	_transactionRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
)

type suiteTransactionService struct {
	suite.Suite
	orderRepositoryMock       *_orderRepositoryMock.OrderRepositoryMock
	transactionRepositoryMock *_transactionRepositoryMock.TransactionRepositoryMock
	orderServiceMock          *_orderServiceMock.OrderServiceMock
	transactionService        TransactionService
}

func newTransactionServiceMock(transaction tr.TransactionRepository, orderRepository or.OrderRepository, orderService os.OrderService) TransactionService {
	return &transactionServiceImpl{
		transactionRepo: transaction,
		orderRepo:       orderRepository,
		orderService:    orderService,
	}
}

func (s *suiteTransactionService) SetupTest() {
	s.orderRepositoryMock = new(_orderRepositoryMock.OrderRepositoryMock)
	s.transactionRepositoryMock = new(_transactionRepositoryMock.TransactionRepositoryMock)
	s.orderServiceMock = new(_orderServiceMock.OrderServiceMock)
	s.transactionService = newTransactionServiceMock(s.transactionRepositoryMock, s.orderRepositoryMock, s.orderServiceMock)
}

func (s *suiteTransactionService) TestCreateTransaction() {
	transactionId := uuid.New().String()
	orderId := uuid.New().String()

	testCase := []struct {
		Name                 string
		Body                 dto.TransactionRequest
		ExpectedErr          error
		FindTransactionErr   error
		CreateTransactionErr error
		UpdateTransactionErr error
		SerOrderStatusErr    error
	}{
		{
			Name: "success",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          nil,
			FindTransactionErr:   customerrors.ErrNotFound,
			CreateTransactionErr: nil,
			UpdateTransactionErr: nil,
			SerOrderStatusErr:    nil,
		},
		{
			Name: "transaction alredy exist, update transaction, and set new order status",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "settlement",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          nil,
			FindTransactionErr:   nil,
			CreateTransactionErr: nil,
			UpdateTransactionErr: nil,
			SerOrderStatusErr:    nil,
		},
		{
			Name: "error when find transaction",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "settlement",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          errors.New("error find transaction"),
			FindTransactionErr:   errors.New("error find transaction"),
			CreateTransactionErr: nil,
			UpdateTransactionErr: nil,
			SerOrderStatusErr:    nil,
		},
		{
			Name: "error when create transaction",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "settlement",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          errors.New("error create transaction"),
			FindTransactionErr:   customerrors.ErrNotFound,
			CreateTransactionErr: errors.New("error create transaction"),
			UpdateTransactionErr: nil,
			SerOrderStatusErr:    nil,
		},
		{
			Name: "error when update transaction",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "settlement",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          errors.New("error update transaction"),
			FindTransactionErr:   nil,
			CreateTransactionErr: nil,
			UpdateTransactionErr: errors.New("error update transaction"),
			SerOrderStatusErr:    nil,
		},
		{
			Name: "error when set new order stratus",
			Body: dto.TransactionRequest{
				TransactionID:     transactionId,
				OrderID:           orderId,
				TransactionStatus: "settlement",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr:          errors.New("error set order status"),
			FindTransactionErr:   nil,
			CreateTransactionErr: nil,
			UpdateTransactionErr: nil,
			SerOrderStatusErr:    errors.New("error set order status"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			mock1 := s.transactionRepositoryMock.On("FindTransaction").Return(v.FindTransactionErr)
			mock2 := s.transactionRepositoryMock.On("CreateTransaction").Return(v.CreateTransactionErr)
			mock3 := s.transactionRepositoryMock.On("UpdateTransaction").Return(v.UpdateTransactionErr)
			mock4 := s.orderServiceMock.On("SetOrderStatus").Return(v.SerOrderStatusErr)

			var ctx context.Context
			err := s.transactionService.CreateTransaction(v.Body, ctx)

			s.Equal(v.ExpectedErr, err)

			mock1.Unset()
			mock2.Unset()
			mock3.Unset()
			mock4.Unset()
		})
	}
}
func (s *suiteTransactionService) TestFindTransaction() {
	UserId := uuid.New()
	TransactionId := uuid.New()
	OrderId := uuid.New()

	testCase := []struct {
		Name                  string
		UserId                string
		ExpectedErr           error
		ExpectedRes           dto.TransactionsResponse
		FindAllTransactionErr error
		FindAllTransactionRes []model.Transaction
	}{
		{
			Name:        "success",
			UserId:      UserId.String(),
			ExpectedErr: nil,
			ExpectedRes: dto.TransactionsResponse{
				{
					TransactionID:     TransactionId.String(),
					PaymentType:       "test",
					TransactionStatus: "pending",
					OrderID:           OrderId.String(),
				},
			},
			FindAllTransactionErr: nil,
			FindAllTransactionRes: []model.Transaction{
				{
					ID:                TransactionId,
					PaymentType:       "test",
					TransactionStatus: "pending",
					OrderID:           OrderId,
				},
			},
		},
		{
			Name:                  "error when find transactions",
			UserId:                UserId.String(),
			ExpectedErr:           errors.New("error find transactions"),
			ExpectedRes:           dto.TransactionsResponse(nil),
			FindAllTransactionErr: errors.New("error find transactions"),
			FindAllTransactionRes: nil,
		},
		{
			Name:                  "transaction nil",
			UserId:                UserId.String(),
			ExpectedErr:           nil,
			ExpectedRes:           nil,
			FindAllTransactionErr: nil,
			FindAllTransactionRes: nil,
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			mock1 := s.transactionRepositoryMock.On("FindAllTransaction").Return(v.FindAllTransactionRes, v.FindAllTransactionErr)

			var ctx context.Context
			res, err := s.transactionService.FindTransaction(v.UserId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			mock1.Unset()
		})
	}
}

func TestSuiteTransactionService(t *testing.T) {
	suite.Run(t, new(suiteTransactionService))
}
