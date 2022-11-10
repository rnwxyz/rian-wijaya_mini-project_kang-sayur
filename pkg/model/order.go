package model

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID `gorm:"primaryKey; type:varchar(50)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UserID        uuid.UUID      `gorm:"; type:varchar(50)"`
	User          User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	CheckpointID  uuid.UUID
	Checkpoint    Checkpoint
	StatusOrderID uint
	StatusOrder   StatusOrder
	ShippingCost  int
	TotalPrice    int
	GrandTotal    int
	OrderDetail   []OrderDetail `gorm:"polymorphic:Order;"`
	Code          string
	Hash          string
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

// gorm hooks or trigger
func (u *Order) AfterCreate(tx *gorm.DB) (err error) {
	for _, ord := range u.OrderDetail { // create order item qty--
		var item Item
		tx.Model(&Item{}).Where("id = ?", ord.ItemID).First(&item)
		newQty := item.Qty - ord.Qty
		tx.Model(&Item{}).Where("id = ?", ord.ItemID).Update("qty", newQty)
	}
	return
}

func (u *Order) AfterUpdate(tx *gorm.DB) (err error) {
	if u.StatusOrderID == 7 { // order cencel item qty++
		var or Order
		tx.Model(&Order{}).Where("id = ?", u.ID).Preload("OrderDetail").First(&or)
		for _, ord := range or.OrderDetail {
			var item Item
			tx.Model(&Item{}).Where("id = ?", ord.ItemID).First(&item)
			newQty := item.Qty + ord.Qty
			tx.Model(&Item{}).Where("id = ?", ord.ItemID).Update("qty", newQty)
		}
		if err != nil {
			panic(err)
		}
		tx.Model(&Order{}).Where("id = ?", u.ID).Update("expired_time", time.Now())
	} else if u.StatusOrderID == 3 { // order ready generate code
		var order Order
		tx.Model(&Order{}).Where("id = ?", u.ID).First(&order)
		// Generating Random string
		ran_str := make([]byte, 5)
		for i := 0; i < 5; i++ {
			ran_str[i] = byte(48 + rand.Intn(10))
		}
		rand := string(ran_str)
		text := order.ID.String() + " " + rand + " " + order.CheckpointID.String()
		str := base64.StdEncoding.EncodeToString([]byte(text))
		if err != nil {
			panic(err)
		}
		tx.Model(&Order{}).Where("id = ?", u.ID).Update("hash", str)
		tx.Model(&Order{}).Where("id = ?", u.ID).Update("code", rand)
		if err != nil {
			panic(err)
		}
		tx.Model(&Order{}).Where("id = ?", u.ID).Update("expired_order", time.Now().Add(12*time.Hour))
	}
	return
}
