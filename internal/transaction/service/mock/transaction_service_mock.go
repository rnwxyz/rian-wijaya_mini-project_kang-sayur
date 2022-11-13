package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/dto"
	"github.com/stretchr/testify/mock"
)

type TransactionServiceMock struct {
	mock.Mock
}

func (b *TransactionServiceMock) CreateTransaction(body dto.TransactionRequest, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *TransactionServiceMock) SetOrderStatus(body dto.TransactionRequest, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *TransactionServiceMock) FindTransaction(id string, ctx context.Context) (dto.TransactionsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.TransactionsResponse), args.Error(1)
}
