package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	importcsv "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type regionServiceImpl struct {
	repo      repository.RegionRepository
	importCsv importcsv.ImportCsv
}

// FindDistrict implements RegionService
func (r *regionServiceImpl) FindDistrict(id *string, ctx context.Context) ([]model.District, error) {
	var district model.District
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, utils.ErrInvalidId
		}
		district.RegencyID = uint(idInt)
	}
	districts, err := r.repo.FindDistrict(&district, ctx)
	if err != nil {
		return nil, err
	}
	return districts, nil
}

// FindProvince implements RegionService
func (r *regionServiceImpl) FindProvince(ctx context.Context) ([]model.Province, error) {
	provinces, err := r.repo.FindProvince(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

// FindRegency implements RegionService
func (r *regionServiceImpl) FindRegency(id *string, ctx context.Context) ([]model.Regency, error) {
	var regency model.Regency
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, utils.ErrInvalidId
		}
		regency.ProvinceID = uint(idInt)
	}
	regencys, err := r.repo.FindRegency(&regency, ctx)
	if err != nil {
		return nil, err
	}
	return regencys, nil
}

// FindVillage implements RegionService
func (r *regionServiceImpl) FindVillage(id *string, ctx context.Context) ([]model.Village, error) {
	var village model.Village
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, utils.ErrInvalidId
		}
		village.DistrictID = uint(idInt)
	}
	villages, err := r.repo.FindVillage(&village, ctx)
	if err != nil {
		return nil, err
	}
	return villages, nil
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
	// limit the placeholder
	const limit = 16200
	placeholders := 0
	for {
		begin := placeholders
		placeholders += limit
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
