package model

// region model
type Province struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Regency struct {
	ID         uint `gorm:"primaryKey"`
	ProvinceID uint
	Province   Province
	Name       string
}

type District struct {
	ID        uint `gorm:"primaryKey"`
	RegencyID uint
	Regency   Regency
	Name      string
}

type Village struct {
	ID         uint `gorm:"primaryKey"`
	DistrictID uint
	District   District
	Name       string
}
