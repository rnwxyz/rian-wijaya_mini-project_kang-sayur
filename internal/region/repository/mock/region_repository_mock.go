package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type RegionRepositoryMock struct {
	mock.Mock
}

func (b *RegionRepositoryMock) CheckIsImported(model interface{}) (bool, error) {
	args := b.Called()
	return args.Bool(0), args.Error(1)
}

func (b *RegionRepositoryMock) ImportRegion(model interface{}) error {
	args := b.Called()
	return args.Error(0)
}

func (b *RegionRepositoryMock) FindProvince(ctx context.Context) ([]model.Province, error) {
	args := b.Called()
	return args.Get(0).([]model.Province), args.Error(1)
}

func (b *RegionRepositoryMock) FindRegency(regency *model.Regency, ctx context.Context) ([]model.Regency, error) {
	args := b.Called()
	return args.Get(0).([]model.Regency), args.Error(1)
}

func (b *RegionRepositoryMock) FindDistrict(district *model.District, ctx context.Context) ([]model.District, error) {
	args := b.Called()
	return args.Get(0).([]model.District), args.Error(1)
}

func (b *RegionRepositoryMock) FindVillage(village *model.Village, ctx context.Context) ([]model.Village, error) {
	args := b.Called()
	return args.Get(0).([]model.Village), args.Error(1)
}
