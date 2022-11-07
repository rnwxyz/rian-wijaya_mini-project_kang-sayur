package model

// region model
type Province struct {
	ID   uint   `gorm:"primaryKey" csv:"province_id" json:"id"`
	Name string `csv:"province_name" json:"name"`
}

type Regency struct {
	ID         uint `gorm:"primaryKey" csv:"regency_id" json:"id"`
	ProvinceID uint `csv:"province_id" json:"province_id"`
	Province   Province
	Name       string `csv:"regency_name" json:"name"`
}

type District struct {
	ID        uint `gorm:"primaryKey" csv:"district_id" json:"id"`
	RegencyID uint `csv:"regency_id" json:"regency_id"`
	Regency   Regency
	Name      string `csv:"district_name" json:"name"`
}

type Village struct {
	ID         uint `gorm:"primaryKey" csv:"village_id" json:"id"`
	DistrictID uint `csv:"district_id" json:"district_id"`
	District   District
	Name       string `csv:"village_name" json:"name"`
}
