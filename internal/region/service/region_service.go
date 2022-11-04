package service

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type RegionService interface {
	importProvince() error
	importRegency() error
	importDistrict() error
	importVillage() error
	FindProvince(id *string, ctx context.Context) ([]model.Province, error)
	FindRegency(id *string, ctx context.Context) ([]model.Regency, error)
	FindDistrict(id *string, ctx context.Context) ([]model.District, error)
	FindVillage(id *string, ctx context.Context) ([]model.Village, error)
}
