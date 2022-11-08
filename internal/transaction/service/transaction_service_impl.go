package service

import (
	"context"

	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	tr "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type transactionServiceImpl struct {
	transactionRepo tr.TransactionRepository
	orderRepo       or.OrderRepository
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
		if err == utils.ErrNotFound {
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
	err = s.SetOrderStatus(body, ctx)
	return err
}

// SetOrderStatus implements TransactionService
func (s *transactionServiceImpl) SetOrderStatus(body dto.TransactionRequest, ctx context.Context) error {
	order := model.Order{
		ID: body.ToModel().OrderID,
	}
	err := s.orderRepo.FindOrderDetail(&order, ctx)
	if err != nil {
		return err
	}
	switch body.TransactionStatus {
	case "capture":
		err := s.orderRepo.OrderWaiting(body.ToModel().OrderID, ctx)
		if err != nil {
			return err
		}
	case "settlement":
		err := s.orderRepo.OrderWaiting(body.ToModel().OrderID, ctx)
		if err != nil {
			return err
		}
	case "deny":
		err := s.orderRepo.CencelOrder(body.ToModel().OrderID, ctx)
		if err != nil {
			return err
		}
	case "cencel":
		err := s.orderRepo.CencelOrder(body.ToModel().OrderID, ctx)
		if err != nil {
			return err
		}
	case "expired":
		err := s.orderRepo.CencelOrder(body.ToModel().OrderID, ctx)
		if err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func NewTransactionService(transaction tr.TransactionRepository, order or.OrderRepository) TransactionService {
	return &transactionServiceImpl{
		transactionRepo: transaction,
		orderRepo:       order,
	}
}
