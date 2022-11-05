package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type districtResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	RegencyID   uint   `json:"regency_id"`
	RegencyName string `json:"regency_name"`
}

func (r *districtResponse) FromModel(model *model.District) {
	r.ID = model.ID
	r.Name = model.Name
	r.RegencyID = model.Regency.ID
	r.RegencyName = model.Regency.Name
}

type DistrictsResponse []districtResponse

func (r *DistrictsResponse) FromModel(model []model.District) {
	for _, each := range model {
		var district districtResponse
		district.FromModel(&each)
		*r = append(*r, district)
	}
}
