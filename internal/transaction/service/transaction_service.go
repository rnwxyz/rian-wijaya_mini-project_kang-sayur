package service

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
)

type TransactionService interface {
	CreateTransaction(body dto.TransactionRequest, ctx context.Context) error
	SetOrderStatus(body dto.TransactionRequest, ctx context.Context) error
}
