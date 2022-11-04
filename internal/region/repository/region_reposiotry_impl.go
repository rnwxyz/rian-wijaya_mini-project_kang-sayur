package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"gorm.io/gorm"
)

type regionRepositoryImpl struct {
	db *gorm.DB
}

// FindDistrict implements RegionRepository
func (r *regionRepositoryImpl) FindDistrict(district *model.District, ctx context.Context) ([]model.District, error) {
	panic("unimplemented")
}

// FindProvince implements RegionRepository
func (r *regionRepositoryImpl) FindProvince(province *model.Province, ctx context.Context) ([]model.Province, error) {
	var provinces []model.Province
	err := r.db.WithContext(ctx).Where(province).Find(provinces).Error
	if err != nil {
		return nil, err
	}
	return provinces, err
}

// FindRegency implements RegionRepository
func (r *regionRepositoryImpl) FindRegency(regency *model.Regency, ctx context.Context) ([]model.Regency, error) {
	panic("unimplemented")
}

// FindVillage implements RegionRepository
func (r *regionRepositoryImpl) FindVillage(village *model.Village, ctx context.Context) ([]model.Village, error) {
	panic("unimplemented")
}

// CheckIsImported implements RegionRepository
func (r *regionRepositoryImpl) CheckIsImported(model interface{}) (bool, error) {
	err := r.db.First(model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// ImportRegion implements RegionRepository
func (r *regionRepositoryImpl) ImportRegion(model interface{}) error {
	err := r.db.Create(model).Error
	return err
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepositoryImpl{
		db: db,
	}
}
