package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Name       string         `gorm:"not null"`
	Email      string         `gorm:"not null;uniqueIndex"`
	Phone      string
	Password   string `gorm:"not null"`
	RoleID     uint
	ProvinceID uint
	Province   Province `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	RegencyID  uint
	Regency    Regency `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	DistrictID uint
	District   District `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	VillageID  uint
	Village    Village `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
