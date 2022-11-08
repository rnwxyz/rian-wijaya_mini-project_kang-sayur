package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type TransactionRepository interface {
	CreateTransaction(transaction *model.Transaction, ctx context.Context) error
	UpdateTransaction(transaction *model.Transaction, ctx context.Context) error
	FindAllTransaction(userId string, ctx context.Context) ([]model.Transaction, error)
	FindTransaction(transaction *model.Transaction, ctx context.Context) error
}
