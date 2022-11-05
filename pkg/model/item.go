package model

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null;unique"`
	CategoryId  uint
	Category    Category
	Description string
	Qty         int
	Price       int
}

type Category struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null;unique"`
	Description string
}
