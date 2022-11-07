package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type regencyResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProvinceID   uint   `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

func (r *regencyResponse) FromModel(model *model.Regency) {
	r.ID = model.ID
	r.Name = model.Name
	r.ProvinceID = model.Province.ID
	r.ProvinceName = model.Province.Name
}

type RegenciesResponse []regencyResponse

func (r *RegenciesResponse) FromModel(model []model.Regency) {
	for _, each := range model {
		var regency regencyResponse
		regency.FromModel(&each)
		*r = append(*r, regency)
	}
}
