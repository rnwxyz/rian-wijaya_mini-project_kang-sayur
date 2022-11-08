package dto

import (
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type TransactionRequest struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime   string `json:"transaction_time"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
	SettlementTime    string `json:"settlement_time"`
}

func (t *TransactionRequest) ToModel() *model.Transaction {
	trId, _ := uuid.Parse(t.TransactionID)
	orId, _ := uuid.Parse(t.OrderID)
	return &model.Transaction{
		ID:                trId,
		OrderID:           orId,
		TransactionStatus: t.TransactionStatus,
		TransactionTime:   t.TransactionTime,
		SignatureKey:      t.SignatureKey,
		PaymentType:       t.PaymentType,
		GrossAmount:       t.GrossAmount,
		SettlementTime:    t.SettlementTime,
	}
}

type TransactionResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime   string `json:"transaction_time"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
	SettlementTime    string `json:"settlement_time"`
}

func (u *TransactionResponse) FromModel(model *model.Transaction) {
	u.TransactionID = model.ID.String()
	u.OrderID = model.OrderID.String()
	u.TransactionStatus = model.TransactionStatus
	u.TransactionTime = model.TransactionTime
	u.SignatureKey = model.SignatureKey
	u.PaymentType = model.PaymentType
	u.GrossAmount = model.GrossAmount
	u.SettlementTime = model.SettlementTime
}

type TransactionsResponse []TransactionResponse

func (u *TransactionsResponse) FromModel(model []model.Transaction) {
	for _, each := range model {
		var transaction TransactionResponse
		transaction.FromModel(&each)
		*u = append(*u, transaction)
	}
}
