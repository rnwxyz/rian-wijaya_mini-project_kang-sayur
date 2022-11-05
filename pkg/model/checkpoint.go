package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Checkpoint struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"not null"`
	Description string
	ProvinceID  uint
	Province    Province
	RegencyID   uint
	Regency     Regency
	DistrictID  uint
	District    District
	VillageID   uint
	Village     Village
	LatLong     string
}
