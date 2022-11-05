package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type villageResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	DistrictID   uint   `json:"district_id"`
	DistrictName string `json:"district_name"`
}

func (r *villageResponse) FromModel(model *model.Village) {
	r.ID = model.ID
	r.Name = model.Name
	r.DistrictID = model.District.ID
	r.DistrictName = model.District.Name
}

type VillagesResponse []villageResponse

func (r *VillagesResponse) FromModel(model []model.Village) {
	for _, each := range model {
		var village villageResponse
		village.FromModel(&each)
		*r = append(*r, village)
	}
}
