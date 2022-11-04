package service

import (
	"context"
	"fmt"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	importcsv "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type regionServiceImpl struct {
	repo      repository.RegionRepository
	importCsv importcsv.ImportCsv
}

// FindDistrict implements RegionService
func (r *regionServiceImpl) FindDistrict(id *string, ctx context.Context) ([]model.District, error) {
	panic("unimplemented")
}

// FindProvince implements RegionService
func (r *regionServiceImpl) FindProvince(id *string, ctx context.Context) ([]model.Province, error) {
	panic("unimplemented")
}

// FindRegency implements RegionService
func (r *regionServiceImpl) FindRegency(id *string, ctx context.Context) ([]model.Regency, error) {
	panic("unimplemented")
}

// FindVillage implements RegionService
func (r *regionServiceImpl) FindVillage(id *string, ctx context.Context) ([]model.Village, error) {
	panic("unimplemented")
}

// importProvince implements RegionService
func (r *regionServiceImpl) importProvince() error {
	isImported, err := r.repo.CheckIsImported(&model.Province{})
	if err != nil {
		return err
	}
	if isImported {
		return nil
	}
	var provinces []model.Province
	err = r.importCsv.UnmarshalCsv(constants.PathProvinceCsv, &provinces)
	if err != nil {
		return err
	}
	err = r.repo.ImportRegion(provinces)
	if err != nil {
		return err
	}
	return nil
}

// importDistrict implements RegionService
func (r *regionServiceImpl) importDistrict() error {
	isImported, err := r.repo.CheckIsImported(&model.District{})
	if err != nil {
		return err
	}
	if isImported {
		return nil
	}
	var districts []model.District
	err = r.importCsv.UnmarshalCsv(constants.PathDistrictCsv, &districts)
	if err != nil {
		return err
	}
	err = r.repo.ImportRegion(districts)
	if err != nil {
		return err
	}
	return nil
}

// importRegency implements RegionService
func (r *regionServiceImpl) importRegency() error {
	isImported, err := r.repo.CheckIsImported(&model.Regency{})
	if err != nil {
		return err
	}
	if isImported {
		return nil
	}
	var regencies []model.Regency
	err = r.importCsv.UnmarshalCsv(constants.PathRegencyCsv, &regencies)
	if err != nil {
		return err
	}
	err = r.repo.ImportRegion(regencies)
	if err != nil {
		return err
	}
	return nil
}

// importVillage implements RegionService
func (r *regionServiceImpl) importVillage() error {
	isImported, err := r.repo.CheckIsImported(&model.Village{})
	if err != nil {
		return err
	}
	if isImported {
		return nil
	}
	var villages []model.Village
	err = r.importCsv.UnmarshalCsv(constants.PathVillageCsv, &villages)
	if err != nil {
		return err
	}
	placeholders := 0
	const safePlaceholder = 16200
	for {
		begin := placeholders
		placeholders += safePlaceholder
		fmt.Println(placeholders)
		if placeholders > len(villages) {
			data := append(villages[0:0], villages[begin:]...)
			err = r.repo.ImportRegion(data)
			if err != nil {
				return err
			}
			return nil
		}
		data := append(villages[0:0], villages[begin:placeholders]...)
		err = r.repo.ImportRegion(data)
		if err != nil {
			return err
		}
	}
}

func NewRegionService(repository repository.RegionRepository, importCsv importcsv.ImportCsv) RegionService {
	NewRegionService := &regionServiceImpl{
		repo:      repository,
		importCsv: importCsv,
	}

	if err := NewRegionService.importProvince(); err != nil {
		panic(err)
	}
	if err := NewRegionService.importRegency(); err != nil {
		panic(err)
	}
	if err := NewRegionService.importDistrict(); err != nil {
		panic(err)
	}
	if err := NewRegionService.importVillage(); err != nil {
		panic(err)
	}

	return NewRegionService
}
