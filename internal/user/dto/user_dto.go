package dto

import (
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type UserSignup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSignup) ToModel() *model.User {
	return &model.User{
		Name:  u.Name,
		Email: u.Email,
	}
}

type UserUpdate struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	ProvinceID uint   `json:"province_id"`
	RegencyID  uint   `json:"regency_id"`
	DistrictID uint   `json:"district_id"`
	VillageID  uint   `json:"village_id"`
}

func (u *UserUpdate) ToModel() *model.User {
	return &model.User{
		Name:       u.Name,
		Phone:      u.Phone,
		ProvinceID: &u.ProvinceID,
		RegencyID:  &u.RegencyID,
		DistrictID: &u.DistrictID,
		VillageID:  &u.VillageID,
	}
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	ProvinceName string    `json:"province_name"`
	RegencyName  string    `json:"regency_name"`
	DistrictName string    `json:"district_name"`
	VillageName  string    `json:"village_name"`
}

func (u *UserResponse) FromModel(model *model.User) {
	u.ID = model.ID
	u.Name = model.Name
	u.Email = model.Email
	u.Phone = model.Phone
	u.ProvinceName = model.Province.Name
	u.RegencyName = model.Regency.Name
	u.DistrictName = model.District.Name
	u.VillageName = model.Village.Name
}

type UsersResponse []UserResponse

func (u *UsersResponse) FromModel(model []model.User) {
	for _, each := range model {
		var user UserResponse
		user.FromModel(&each)
		*u = append(*u, user)
	}
}
