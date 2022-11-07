package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                uuid.UUID `gorm:"primaryKey; type:varchar(50)"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	OrderID           uuid.UUID      `gorm:"; type:varchar(50)"`
	Order             Order
	TransactionStatus string
	TransactionTime   string
	SignatureKey      string
	PaymentType       string
	GrossAmount       string
	SettlementTime    string
}
