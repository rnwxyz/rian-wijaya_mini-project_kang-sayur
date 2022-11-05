package dto

import (
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type CheckpointRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ProvinceID  uint   `json:"province_id" validate:"required"`
	RegencyID   uint   `json:"regency_id" validate:"required"`
	DistrictID  uint   `json:"district_id" validate:"required"`
	VillageID   uint   `json:"village_id" validate:"required"`
	LatLong     string `json:"lat_long" validate:"required"`
}

func (u *CheckpointRequest) ToModel() *model.Checkpoint {
	return &model.Checkpoint{
		Name:        u.Name,
		Description: u.Description,
		ProvinceID:  u.ProvinceID,
		RegencyID:   u.RegencyID,
		DistrictID:  u.DistrictID,
		VillageID:   u.VillageID,
		LatLong:     u.LatLong,
	}
}

type CheckpointResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ProvinceName string    `json:"province_name"`
	RegencyName  string    `json:"regency_name"`
	DistrictName string    `json:"district_name"`
	VillageName  string    `json:"village_name"`
	LatLong      string    `json:"lat_long"`
}

func (u *CheckpointResponse) FromModel(model *model.Checkpoint) {
	u.ID = model.ID
	u.Name = model.Name
	u.Description = model.Description
	u.ProvinceName = model.Province.Name
	u.RegencyName = model.Regency.Name
	u.DistrictName = model.District.Name
	u.VillageName = model.Village.Name
	u.LatLong = model.LatLong
}

type CheckpointsResponse []CheckpointResponse

func (u *CheckpointsResponse) FromModel(model []model.Checkpoint) {
	for _, each := range model {
		var checkpoint CheckpointResponse
		checkpoint.FromModel(&each)
		*u = append(*u, checkpoint)
	}
}
