package repository

import (
	"context"
	"strings"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

// CreateTransaction implements TransactionRepository
func (r *transactionRepositoryImpl) CreateTransaction(transaction *model.Transaction, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(transaction).Error
	if err != nil {
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return utils.ErrBadRequestBody
		}
		return err
	}
	return nil
}

// CreateTransaction implements TransactionRepository
func (r *transactionRepositoryImpl) UpdateTransaction(transaction *model.Transaction, ctx context.Context) error {
	err := r.db.WithContext(ctx).Save(transaction).Error
	if err != nil {
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return utils.ErrBadRequestBody
		}
		return err
	}
	return nil
}

// FindAllTransaction implements TransactionRepository
func (r *transactionRepositoryImpl) FindAllTransaction(ctx context.Context) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.WithContext(ctx).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// FindTransaction implements TransactionRepository
func (r *transactionRepositoryImpl) FindTransaction(transaction *model.Transaction, ctx context.Context) error {
	err := r.db.WithContext(ctx).First(transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}
