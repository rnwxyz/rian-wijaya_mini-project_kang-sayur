package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID `gorm:"primaryKey; type:varchar(50)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UserID        uuid.UUID      `gorm:"; type:varchar(50)"`
	User          User           `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
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
	OrderType string
	Order     Order
	ItemID    uint
	Item      Item `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	Qty       int
	Price     int
	Total     int
}

func (u *Order) AfterCreate(tx *gorm.DB) (err error) {
	for _, ord := range u.OrderDetail {
		var item Item
		tx.Model(&Item{}).Where("id = ?", ord.ItemID).First(&item)
		newQty := item.Qty - ord.Qty
		tx.Model(&Item{}).Where("id = ?", ord.ItemID).Update("qty", newQty)
	}
	return
}

func (u *Order) AfterUpdate(tx *gorm.DB) (err error) {
	if u.StatusOrderID == constants.Cencel_status_order_id {
		var or Order
		tx.Model(&Order{}).Where("id = ?", u.ID).Preload("OrderDetail").First(&or)
		for _, ord := range or.OrderDetail {
			var item Item
			tx.Model(&Item{}).Where("id = ?", ord.ItemID).First(&item)
			newQty := item.Qty + ord.Qty
			tx.Model(&Item{}).Where("id = ?", ord.ItemID).Update("qty", newQty)
		}
	}
	return
}
