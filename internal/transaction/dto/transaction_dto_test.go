package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRequest_ToModel(t *testing.T) {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name     string
		Dto      TransactionRequest
		Expected *model.Transaction
	}{
		{
			Name: "all filled",
			Dto: TransactionRequest{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			},
			Expected: &model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			},
		},
		{
			Name: "some filled",
			Dto: TransactionRequest{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
			},
			Expected: &model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
			},
		},
		{
			Name:     "dto empty",
			Dto:      TransactionRequest{},
			Expected: &model.Transaction{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			result := v.Dto.ToModel()
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestTransactionResponse_FromModel(t *testing.T) {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name     string
		Model    *model.Transaction
		Expected TransactionResponse
	}{
		{
			Name: "all filled",
			Expected: TransactionResponse{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			},
			Model: &model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			},
		},
		{
			Name: "some filled",
			Expected: TransactionResponse{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
			},
			Model: &model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
			},
		},
		{
			Name: "model empty",
			Expected: TransactionResponse{
				TransactionID: "00000000-0000-0000-0000-000000000000",
				OrderID:       "00000000-0000-0000-0000-000000000000",
			},
			Model: &model.Transaction{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result TransactionResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestTransactionsResponse_FromModel(t *testing.T) {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name     string
		Model    []model.Transaction
		Expected TransactionsResponse
	}{
		{
			Name: "all filled",
			Expected: TransactionsResponse{{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			}},
			Model: []model.Transaction{{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
				TransactionTime:   "2022-01-01 00:00:00",
				SignatureKey:      "003592974b7cb5956dbfe36c000b",
				PaymentType:       "qris",
				GrossAmount:       "10000",
				SettlementTime:    "",
			}},
		},
		{
			Name: "some filled",
			Model: []model.Transaction{{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "pending",
			}},
			Expected: TransactionsResponse{{
				TransactionID:     transactionId.String(),
				OrderID:           orderId.String(),
				TransactionStatus: "pending",
			}},
		},
		{
			Name:     "model empty",
			Model:    []model.Transaction{},
			Expected: TransactionsResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result TransactionsResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
