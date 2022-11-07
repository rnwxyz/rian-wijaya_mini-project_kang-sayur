package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey; type:varchar(50)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Name       string         `gorm:"not null"`
	Email      string         `gorm:"not null;unique"`
	Phone      string
	Password   string `gorm:"not null"`
	RoleID     uint
	Role       Role
	ProvinceID *uint
	Province   Province
	RegencyID  *uint
	Regency    Regency
	DistrictID *uint
	District   District
	VillageID  *uint
	Village    Village
}

type Role struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Description string
}
