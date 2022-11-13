package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type RegionServiceMock struct {
	mock.Mock
}

func (b *RegionServiceMock) ImportProvince() error {
	args := b.Called()
	return args.Error(0)
}

func (b *RegionServiceMock) ImportRegency() error {
	args := b.Called()
	return args.Error(0)
}

func (b *RegionServiceMock) ImportDistrict() error {
	args := b.Called()
	return args.Error(0)
}

func (b *RegionServiceMock) ImportVillage() error {
	args := b.Called()
	return args.Error(0)
}

func (b *RegionServiceMock) FindProvince(ctx context.Context) ([]model.Province, error) {
	args := b.Called()
	return args.Get(0).([]model.Province), args.Error(1)
}

func (b *RegionServiceMock) FindRegency(id *string, ctx context.Context) (dto.RegenciesResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.RegenciesResponse), args.Error(1)
}

func (b *RegionServiceMock) FindDistrict(id *string, ctx context.Context) (dto.DistrictsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.DistrictsResponse), args.Error(1)
}

func (b *RegionServiceMock) FindVillage(id *string, ctx context.Context) (dto.VillagesResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.VillagesResponse), args.Error(1)
}
