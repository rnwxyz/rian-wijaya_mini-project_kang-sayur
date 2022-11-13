package service

import (
	"context"
	"strconv"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	importcsv "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type regionServiceImpl struct {
	repo      repository.RegionRepository
	importCsv importcsv.ImportCsv
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
func (r *regionServiceImpl) FindRegency(id *string, ctx context.Context) (dto.RegenciesResponse, error) {
	var regency model.Regency
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, customerrors.ErrInvalidId
		}
		regency.ProvinceID = uint(idInt)
	}
	regenciesModel, err := r.repo.FindRegency(&regency, ctx)
	if err != nil {
		return nil, err
	}
	var regencies dto.RegenciesResponse
	regencies.FromModel(regenciesModel)
	return regencies, nil
}

// FindDistrict implements RegionService
func (r *regionServiceImpl) FindDistrict(id *string, ctx context.Context) (dto.DistrictsResponse, error) {
	var district model.District
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, customerrors.ErrInvalidId
		}
		district.RegencyID = uint(idInt)
	}
	districtsModel, err := r.repo.FindDistrict(&district, ctx)
	if err != nil {
		return nil, err
	}
	var districts dto.DistrictsResponse
	districts.FromModel(districtsModel)
	return districts, nil
}

// FindVillage implements RegionService
func (r *regionServiceImpl) FindVillage(id *string, ctx context.Context) (dto.VillagesResponse, error) {
	var village model.Village
	if id != nil {
		idInt, err := strconv.Atoi(*id)
		if err != nil {
			return nil, customerrors.ErrInvalidId
		}
		village.DistrictID = uint(idInt)
	}
	villagesModel, err := r.repo.FindVillage(&village, ctx)
	if err != nil {
		return nil, err
	}
	var villages dto.VillagesResponse
	villages.FromModel(villagesModel)
	return villages, nil
}

// ImportProvince implements RegionService
func (r *regionServiceImpl) ImportProvince() error {
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

// ImportDistrict implements RegionService
func (r *regionServiceImpl) ImportDistrict() error {
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

// ImportRegency implements RegionService
func (r *regionServiceImpl) ImportRegency() error {
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

// ImportVillage implements RegionService
func (r *regionServiceImpl) ImportVillage() error {
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

	if err := NewRegionService.ImportProvince(); err != nil {
		panic(err)
	}
	if err := NewRegionService.ImportRegency(); err != nil {
		panic(err)
	}
	if err := NewRegionService.ImportDistrict(); err != nil {
		panic(err)
	}
	if err := NewRegionService.ImportVillage(); err != nil {
		panic(err)
	}

	return NewRegionService
}
