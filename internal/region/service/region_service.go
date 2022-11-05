package service

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type RegionService interface {
	importProvince() error
	importRegency() error
	importDistrict() error
	importVillage() error
	FindProvince(ctx context.Context) ([]model.Province, error)
	FindRegency(id *string, ctx context.Context) (dto.RegenciesResponse, error)
	FindDistrict(id *string, ctx context.Context) (dto.DistrictsResponse, error)
	FindVillage(id *string, ctx context.Context) (dto.VillagesResponse, error)
}
