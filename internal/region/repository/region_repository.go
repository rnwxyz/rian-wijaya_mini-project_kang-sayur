package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type RegionRepository interface {
	CheckIsImported(model interface{}) (bool, error)
	ImportRegion(model interface{}) error
	FindProvince(province *model.Province, ctx context.Context) ([]model.Province, error)
	FindRegency(regency *model.Regency, ctx context.Context) ([]model.Regency, error)
	FindDistrict(district *model.District, ctx context.Context) ([]model.District, error)
	FindVillage(village *model.Village, ctx context.Context) ([]model.Village, error)
}
