package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (b *TransactionRepositoryMock) CreateTransaction(transaction *model.Transaction, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}

func (b *TransactionRepositoryMock) UpdateTransaction(transaction *model.Transaction, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}

func (b *TransactionRepositoryMock) FindAllTransaction(userId string, ctx context.Context) ([]model.Transaction, error) {
	args := b.Called()

	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (b *TransactionRepositoryMock) FindTransaction(transaction *model.Transaction, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}
