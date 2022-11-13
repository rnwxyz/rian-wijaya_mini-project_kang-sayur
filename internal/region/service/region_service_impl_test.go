package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository"
	regionRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository/mock"
	importcsv "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv"
	importcsvMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
)

type suiteRegionService struct {
	suite.Suite
	regionRepositoryMock *regionRepositoryMock.RegionRepositoryMock
	importcsvMock        *importcsvMock.ImportCsvMock
	regionService        RegionService
}

func newRegionServiceMock(regionRepository repository.RegionRepository, importCsv importcsv.ImportCsv) RegionService {
	return &regionServiceImpl{
		repo:      regionRepository,
		importCsv: importCsv,
	}
}

func (s *suiteRegionService) SetupSuit() {
	s.regionRepositoryMock = new(regionRepositoryMock.RegionRepositoryMock)
	s.importcsvMock = new(importcsvMock.ImportCsvMock)
	s.regionService = newRegionServiceMock(s.regionRepositoryMock, s.importcsvMock)
}

func (s *suiteRegionService) TearDown() {
	s.regionRepositoryMock = nil
	s.importcsvMock = nil
	s.regionService = nil
}

func (s *suiteRegionService) TestImportProvince() {
	testCase := []struct {
		Name               string
		ExpectedErr        error
		CheckIsImportedErr error
		CheckIsImportedRes bool
		UnmarshalCsvErr    error
		ImportRegionErr    error
	}{
		{
			Name:               "import province",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "province is imported",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is imported",
			ExpectedErr:        errors.New("error check is imported"),
			CheckIsImportedErr: errors.New("error check is imported"),
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is unmarshal CSV file",
			ExpectedErr:        errors.New("error unmarshal csv"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    errors.New("error unmarshal csv"),
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when insert province",
			ExpectedErr:        errors.New("error insert province"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    errors.New("error insert province"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("CheckIsImported").Return(v.CheckIsImportedRes, v.CheckIsImportedErr)
			s.importcsvMock.On("UnmarshalCsv").Return(v.UnmarshalCsvErr)
			s.regionRepositoryMock.On("ImportRegion").Return(v.ImportRegionErr)

			err := s.regionService.ImportProvince()

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteRegionService) TestImportRegency() {
	testCase := []struct {
		Name               string
		ExpectedErr        error
		CheckIsImportedErr error
		CheckIsImportedRes bool
		UnmarshalCsvErr    error
		ImportRegionErr    error
	}{
		{
			Name:               "import regency",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "regency is imported",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is imported",
			ExpectedErr:        errors.New("error check is imported"),
			CheckIsImportedErr: errors.New("error check is imported"),
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is unmarshal CSV file",
			ExpectedErr:        errors.New("error unmarshal csv"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    errors.New("error unmarshal csv"),
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when insert regency",
			ExpectedErr:        errors.New("error insert regency"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    errors.New("error insert regency"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("CheckIsImported").Return(v.CheckIsImportedRes, v.CheckIsImportedErr)
			s.importcsvMock.On("UnmarshalCsv").Return(v.UnmarshalCsvErr)
			s.regionRepositoryMock.On("ImportRegion").Return(v.ImportRegionErr)

			err := s.regionService.ImportRegency()

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteRegionService) TestImportDistrict() {
	testCase := []struct {
		Name               string
		ExpectedErr        error
		CheckIsImportedErr error
		CheckIsImportedRes bool
		UnmarshalCsvErr    error
		ImportRegionErr    error
	}{
		{
			Name:               "import district",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "district is imported",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is imported",
			ExpectedErr:        errors.New("error check is imported"),
			CheckIsImportedErr: errors.New("error check is imported"),
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is unmarshal CSV file",
			ExpectedErr:        errors.New("error unmarshal csv"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    errors.New("error unmarshal csv"),
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when insert district",
			ExpectedErr:        errors.New("error insert district"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    errors.New("error insert district"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("CheckIsImported").Return(v.CheckIsImportedRes, v.CheckIsImportedErr)
			s.importcsvMock.On("UnmarshalCsv").Return(v.UnmarshalCsvErr)
			s.regionRepositoryMock.On("ImportRegion").Return(v.ImportRegionErr)

			err := s.regionService.ImportDistrict()

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteRegionService) TestImportVillage() {
	testCase := []struct {
		Name               string
		ExpectedErr        error
		CheckIsImportedErr error
		CheckIsImportedRes bool
		UnmarshalCsvErr    error
		ImportRegionErr    error
	}{
		{
			Name:               "import village",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "village is imported",
			ExpectedErr:        nil,
			CheckIsImportedErr: nil,
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is imported",
			ExpectedErr:        errors.New("error check is imported"),
			CheckIsImportedErr: errors.New("error check is imported"),
			CheckIsImportedRes: true,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when check is unmarshal CSV file",
			ExpectedErr:        errors.New("error unmarshal csv"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    errors.New("error unmarshal csv"),
			ImportRegionErr:    nil,
		},
		{
			Name:               "error when insert village",
			ExpectedErr:        errors.New("error insert village"),
			CheckIsImportedErr: nil,
			CheckIsImportedRes: false,
			UnmarshalCsvErr:    nil,
			ImportRegionErr:    errors.New("error insert village"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("CheckIsImported").Return(v.CheckIsImportedRes, v.CheckIsImportedErr)
			s.importcsvMock.On("UnmarshalCsv").Return(v.UnmarshalCsvErr)
			s.regionRepositoryMock.On("ImportRegion").Return(v.ImportRegionErr)

			err := s.regionService.ImportVillage()

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}

func (s *suiteRegionService) TestFindProvince() {
	testCase := []struct {
		Name            string
		ExpectedErr     error
		ExpectedRes     []model.Province
		FindProvinceErr error
		FindProvinceRes []model.Province
	}{
		{
			Name:        "find province success",
			ExpectedErr: nil,
			ExpectedRes: []model.Province{
				{
					ID:   1,
					Name: "provinceA",
				},
				{
					ID:   2,
					Name: "provinceB",
				},
			},
			FindProvinceErr: nil,
			FindProvinceRes: []model.Province{
				{
					ID:   1,
					Name: "provinceA",
				},
				{
					ID:   2,
					Name: "provinceB",
				},
			},
		},
		{
			Name:            "error when find province",
			ExpectedErr:     errors.New("error when find province"),
			ExpectedRes:     []model.Province(nil),
			FindProvinceErr: errors.New("error when find province"),
			FindProvinceRes: []model.Province{},
		},
		{
			Name:            "empty province",
			ExpectedErr:     nil,
			ExpectedRes:     []model.Province{},
			FindProvinceErr: nil,
			FindProvinceRes: []model.Province{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("FindProvince").Return(v.FindProvinceRes, v.FindProvinceErr)

			var ctx context.Context
			res, err := s.regionService.FindProvince(ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteRegionService) TestFindRegency() {
	testCase := []struct {
		Name           string
		ExpectedErr    error
		ProvinceId     string
		ExpectedRes    dto.RegenciesResponse
		FindRegencyErr error
		FindRegencyRes []model.Regency
	}{
		{
			Name:        "find regency success",
			ExpectedErr: nil,
			ProvinceId:  "1",
			ExpectedRes: dto.RegenciesResponse{
				{
					ID:           1,
					Name:         "regencyA",
					ProvinceID:   1,
					ProvinceName: "provinceA",
				},
				{
					ID:           2,
					Name:         "regencyB",
					ProvinceID:   1,
					ProvinceName: "provinceA",
				},
			},
			FindRegencyErr: nil,
			FindRegencyRes: []model.Regency{
				{
					ID:         1,
					Name:       "regencyA",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "provinceA",
					},
				},
				{
					ID:         2,
					Name:       "regencyB",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "provinceA",
					},
				},
			},
		},
		{
			Name:           "error when find Regency",
			ProvinceId:     "1",
			ExpectedErr:    errors.New("error when find Regency"),
			ExpectedRes:    dto.RegenciesResponse(nil),
			FindRegencyErr: errors.New("error when find Regency"),
			FindRegencyRes: []model.Regency{},
		},
		{
			Name:           "empty Regency",
			ProvinceId:     "1",
			ExpectedErr:    nil,
			ExpectedRes:    dto.RegenciesResponse(nil),
			FindRegencyErr: nil,
			FindRegencyRes: []model.Regency{},
		},
		{
			Name:           "invalid id",
			ProvinceId:     "a",
			ExpectedErr:    customerrors.ErrInvalidId,
			ExpectedRes:    dto.RegenciesResponse(nil),
			FindRegencyErr: nil,
			FindRegencyRes: []model.Regency{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("FindRegency").Return(v.FindRegencyRes, v.FindRegencyErr)

			var ctx context.Context
			res, err := s.regionService.FindRegency(&v.ProvinceId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func (s *suiteRegionService) TestFindDistrict() {
	testCase := []struct {
		Name            string
		ExpectedErr     error
		RegencyId       string
		ExpectedRes     dto.DistrictsResponse
		FindDistrictErr error
		FindDistrictRes []model.District
	}{
		{
			Name:        "find district success",
			ExpectedErr: nil,
			RegencyId:   "1",
			ExpectedRes: dto.DistrictsResponse{
				{
					ID:          1,
					Name:        "districtA",
					RegencyID:   1,
					RegencyName: "regencyA",
				},
				{
					ID:          2,
					Name:        "districtB",
					RegencyID:   1,
					RegencyName: "regencyA",
				},
			},
			FindDistrictErr: nil,
			FindDistrictRes: []model.District{
				{
					ID:        1,
					Name:      "districtA",
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regencyA",
					},
				},
				{
					ID:        2,
					Name:      "districtB",
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regencyA",
					},
				},
			},
		},
		{
			Name:            "error when find district",
			RegencyId:       "1",
			ExpectedErr:     errors.New("error when find district"),
			ExpectedRes:     dto.DistrictsResponse(nil),
			FindDistrictErr: errors.New("error when find district"),
			FindDistrictRes: []model.District{},
		},
		{
			Name:            "empty district",
			RegencyId:       "1",
			ExpectedErr:     nil,
			ExpectedRes:     dto.DistrictsResponse(nil),
			FindDistrictErr: nil,
			FindDistrictRes: []model.District{},
		},
		{
			Name:            "invalid id",
			RegencyId:       "a",
			ExpectedErr:     customerrors.ErrInvalidId,
			ExpectedRes:     dto.DistrictsResponse(nil),
			FindDistrictErr: nil,
			FindDistrictRes: []model.District{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("FindDistrict").Return(v.FindDistrictRes, v.FindDistrictErr)

			var ctx context.Context
			res, err := s.regionService.FindDistrict(&v.RegencyId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func (s *suiteRegionService) TestFindVillage() {
	testCase := []struct {
		Name           string
		ExpectedErr    error
		DistrictId     string
		ExpectedRes    dto.VillagesResponse
		FindVillageErr error
		FindVillageRes []model.Village
	}{
		{
			Name:        "find village success",
			ExpectedErr: nil,
			DistrictId:  "1",
			ExpectedRes: dto.VillagesResponse{
				{
					ID:           1,
					Name:         "villageA",
					DistrictID:   1,
					DistrictName: "districtA",
				},
				{
					ID:           2,
					Name:         "villageB",
					DistrictID:   1,
					DistrictName: "districtA",
				},
			},
			FindVillageErr: nil,
			FindVillageRes: []model.Village{
				{
					ID:         1,
					Name:       "villageA",
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "districtA",
					},
				},
				{
					ID:         2,
					Name:       "villageB",
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "districtA",
					},
				},
			},
		},
		{
			Name:           "error when find village",
			DistrictId:     "1",
			ExpectedErr:    errors.New("error when find village"),
			ExpectedRes:    dto.VillagesResponse(nil),
			FindVillageErr: errors.New("error when find village"),
			FindVillageRes: []model.Village{},
		},
		{
			Name:           "empty village",
			DistrictId:     "1",
			ExpectedErr:    nil,
			ExpectedRes:    dto.VillagesResponse(nil),
			FindVillageErr: nil,
			FindVillageRes: []model.Village{},
		},
		{
			Name:           "invalid id",
			DistrictId:     "a",
			ExpectedErr:    customerrors.ErrInvalidId,
			ExpectedRes:    dto.VillagesResponse(nil),
			FindVillageErr: nil,
			FindVillageRes: []model.Village{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.regionRepositoryMock.On("FindVillage").Return(v.FindVillageRes, v.FindVillageErr)

			var ctx context.Context
			res, err := s.regionService.FindVillage(&v.DistrictId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestRegionService(t *testing.T) {
	suite.Run(t, new(suiteRegionService))
}
