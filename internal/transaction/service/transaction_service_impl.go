package service

import (
	"context"

	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	os "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	tr "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/repository"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type transactionServiceImpl struct {
	transactionRepo tr.TransactionRepository
	orderRepo       or.OrderRepository
	orderService    os.OrderService
}

// FindTransaction implements TransactionService
func (s *transactionServiceImpl) FindTransaction(id string, ctx context.Context) (dto.TransactionsResponse, error) {
	transactions, err := s.transactionRepo.FindAllTransaction(id, ctx)
	if err != nil {
		return nil, err
	}
	var transactionsResponse dto.TransactionsResponse
	transactionsResponse.FromModel(transactions)
	return transactionsResponse, nil
}

// CreateTransaction implements TransactionService
func (s *transactionServiceImpl) CreateTransaction(body dto.TransactionRequest, ctx context.Context) error {
	model := body.ToModel()
	err := s.transactionRepo.FindTransaction(model, ctx)
	if err != nil {
		if err == customerrors.ErrNotFound {
			err = s.transactionRepo.CreateTransaction(model, ctx)
			return err
		}
		return err
	}
	makeTransaction := body.ToModel()
	err = s.transactionRepo.UpdateTransaction(makeTransaction, ctx)
	if err != nil {
		return err
	}
	err = s.orderService.SetOrderStatus(makeTransaction.OrderID, makeTransaction.TransactionStatus, ctx)
	return err
}

func NewTransactionService(transaction tr.TransactionRepository, orderRepository or.OrderRepository, orderService os.OrderService) TransactionService {
	return &transactionServiceImpl{
		transactionRepo: transaction,
		orderRepo:       orderRepository,
		orderService:    orderService,
	}
}
