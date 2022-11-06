package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UserID        uuid.UUID
	User          User `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	StatusOrderID uint
	StatusOrder   StatusOrder
	ShippingCost  int
	TotalPrice    int
	GrandTotal    int
	OrderDetail   []OrderDetail `gorm:"polymorphic:Order;"`
	Code          string
	ExpiredOrder  time.Time
}

type StatusOrder struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null;unique"`
	Description string
}

type OrderDetail struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	OrderID   uuid.UUID
	Order     Order
	ItemID    uint
	Item      Item `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	Qty       int
	Price     int
	Total     int
}
