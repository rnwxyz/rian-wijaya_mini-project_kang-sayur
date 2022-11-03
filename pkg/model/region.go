package model

// region model
type Province struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Regency struct {
	ID         uint `gorm:"primaryKey"`
	ProvinceID uint
	Province   Province `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name       string
}

type District struct {
	ID        uint `gorm:"primaryKey"`
	RegencyID uint
	Regency   Regency `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name      string
}

type Village struct {
	ID         uint `gorm:"primaryKey"`
	DistrictID uint
	District   District `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name       string
}
