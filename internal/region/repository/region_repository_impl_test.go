package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteRegionRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *regionRepositoryImpl
}

func (s *suiteRegionRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &regionRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteRegionRepository) TearDown() {
	s.mock = nil
	s.repository = nil
}

func (s *suiteRegionRepository) TestImportRegion() {
	testCase := []struct {
		Name        string
		Body        interface{}
		ExpectedErr error
		MockReturn  error
	}{
		//example use model province
		{
			Name: "import success",
			Body: []model.Province{
				{
					ID:   1,
					Name: "provinceA",
				},
				{
					ID:   2,
					Name: "provinceB",
				},
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "import error",
			Body: []model.Province{
				{
					ID:   1,
					Name: "provinceA",
				},
				{
					ID:   2,
					Name: "provinceB",
				},
			},
			ExpectedErr: errors.New("error"),
			MockReturn:  errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `provinces` (`name`,`id`) VALUES (?,?),(?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}
			err := s.repository.ImportRegion(v.Body)

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}

func (s *suiteRegionRepository) TestFindProvince() {
	testCase := []struct {
		Name        string
		ExpectedErr error
		ExpectedRes []model.Province
		MockErr     error
		MockRes     *sqlmock.Rows
	}{
		{
			Name:        "success",
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
			MockErr: nil,
			MockRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "provinceA").
				AddRow(2, "provinceB"),
		},
		{
			Name:        "nil province",
			ExpectedErr: nil,
			ExpectedRes: []model.Province{},
			MockErr:     nil,
			MockRes:     sqlmock.NewRows(nil),
		},
		{
			Name:        "error",
			ExpectedErr: errors.New("error"),
			ExpectedRes: []model.Province(nil),
			MockErr:     errors.New("error"),
			MockRes:     sqlmock.NewRows(nil),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces`"))
			db.WillReturnRows(v.MockRes)
			db.WillReturnError(v.MockErr)
			var ctx context.Context
			res, err := s.repository.FindProvince(ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteRegionRepository) TestFindRegency() {
	testCase := []struct {
		Name                  string
		ExpectedErr           error
		Body                  model.Regency
		ExpectedRes           []model.Regency
		MockSelectRegencyErr  error
		MockSelectRegencyRes  *sqlmock.Rows
		MockSelectProvinceErr error
		MockSelectProvinceRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			Body: model.Regency{
				ProvinceID: 1,
			},
			ExpectedRes: []model.Regency{
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
			MockSelectRegencyErr: nil,
			MockSelectRegencyRes: sqlmock.NewRows([]string{"id", "name", "province_id"}).
				AddRow(1, "regencyA", 1).
				AddRow(2, "regencyB", 1),
			MockSelectProvinceErr: nil,
			MockSelectProvinceRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "provinceA"),
		},
		{
			Name:        "nil regency",
			ExpectedErr: nil,
			Body: model.Regency{
				ProvinceID: 1,
			},
			ExpectedRes:           []model.Regency{},
			MockSelectRegencyErr:  nil,
			MockSelectRegencyRes:  sqlmock.NewRows([]string{"id", "name", "province_id"}),
			MockSelectProvinceErr: nil,
			MockSelectProvinceRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "provinceA"),
		},
		{
			Name:        "error",
			ExpectedErr: errors.New("error"),
			Body: model.Regency{
				ProvinceID: 1,
			},
			ExpectedRes:           []model.Regency(nil),
			MockSelectRegencyErr:  errors.New("error"),
			MockSelectRegencyRes:  sqlmock.NewRows([]string{"id", "name", "province_id"}),
			MockSelectProvinceErr: nil,
			MockSelectProvinceRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "provinceA"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			selectRegency := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`province_id` = ?"))
			selectRegency.WillReturnRows(v.MockSelectRegencyRes)
			selectRegency.WillReturnError(v.MockSelectRegencyErr)

			selectProvince := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?"))
			selectProvince.WillReturnRows(v.MockSelectProvinceRes)
			selectProvince.WillReturnError(v.MockSelectProvinceErr)

			var ctx context.Context
			res, err := s.repository.FindRegency(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteRegionRepository) TestFindDistrict() {
	testCase := []struct {
		Name                  string
		ExpectedErr           error
		Body                  model.District
		ExpectedRes           []model.District
		MockSelectDistrictErr error
		MockSelectDistrictRes *sqlmock.Rows
		MockSelectRegencyErr  error
		MockSelectRegencyRes  *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			Body: model.District{
				RegencyID: 1,
			},
			ExpectedRes: []model.District{
				{
					ID:        1,
					Name:      "districtA",
					RegencyID: 1,
					Regency: model.Regency{
						ID:         1,
						Name:       "regencyA",
						ProvinceID: 1,
					},
				},
				{
					ID:        2,
					Name:      "districtB",
					RegencyID: 1,
					Regency: model.Regency{
						ID:         1,
						Name:       "regencyA",
						ProvinceID: 1,
					},
				},
			},
			MockSelectDistrictErr: nil,
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}).
				AddRow(1, "districtA", 1).
				AddRow(2, "districtB", 1),
			MockSelectRegencyErr: nil,
			MockSelectRegencyRes: sqlmock.NewRows([]string{"id", "name", "province_id"}).
				AddRow(1, "regencyA", 1),
		},
		{
			Name:        "nil District",
			ExpectedErr: nil,
			Body: model.District{
				RegencyID: 1,
			},
			ExpectedRes:           []model.District{},
			MockSelectDistrictErr: nil,
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}),
			MockSelectRegencyErr:  nil,
			MockSelectRegencyRes: sqlmock.NewRows([]string{"id", "name", "province_id"}).
				AddRow(1, "regencyA", 1),
		},
		{
			Name:        "error",
			ExpectedErr: errors.New("error"),
			Body: model.District{
				RegencyID: 1,
			},
			ExpectedRes:           []model.District(nil),
			MockSelectDistrictErr: errors.New("error"),
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}),
			MockSelectRegencyErr:  nil,
			MockSelectRegencyRes: sqlmock.NewRows([]string{"id", "name", "province_id"}).
				AddRow(1, "regencyA", 1),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			selectDistrict := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`regency_id` = ?"))
			selectDistrict.WillReturnRows(v.MockSelectDistrictRes)
			selectDistrict.WillReturnError(v.MockSelectDistrictErr)

			selectRegency := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?"))
			selectRegency.WillReturnRows(v.MockSelectRegencyRes)
			selectRegency.WillReturnError(v.MockSelectRegencyErr)

			var ctx context.Context
			res, err := s.repository.FindDistrict(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteRegionRepository) TestFindVillage() {
	testCase := []struct {
		Name                  string
		ExpectedErr           error
		Body                  model.Village
		ExpectedRes           []model.Village
		MockSelectVillageErr  error
		MockSelectVillageRes  *sqlmock.Rows
		MockSelectDistrictErr error
		MockSelectDistrictRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			Body: model.Village{
				DistrictID: 1,
			},
			ExpectedRes: []model.Village{
				{
					ID:         1,
					Name:       "villageA",
					DistrictID: 1,
					District: model.District{
						ID:        1,
						Name:      "districtA",
						RegencyID: 1,
					},
				},
				{
					ID:         2,
					Name:       "villageB",
					DistrictID: 1,
					District: model.District{
						ID:        1,
						Name:      "districtA",
						RegencyID: 1,
					},
				},
			},
			MockSelectVillageErr: nil,
			MockSelectVillageRes: sqlmock.NewRows([]string{"id", "name", "district_id"}).
				AddRow(1, "villageA", 1).
				AddRow(2, "villageB", 1),
			MockSelectDistrictErr: nil,
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}).
				AddRow(1, "districtA", 1),
		},
		{
			Name:        "nil village",
			ExpectedErr: nil,
			Body: model.Village{
				DistrictID: 1,
			},
			ExpectedRes:           []model.Village{},
			MockSelectVillageErr:  nil,
			MockSelectVillageRes:  sqlmock.NewRows([]string{"id", "name", "district_id"}),
			MockSelectDistrictErr: nil,
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}).
				AddRow(1, "districtA", 1),
		},
		{
			Name:        "error",
			ExpectedErr: errors.New("error"),
			Body: model.Village{
				DistrictID: 1,
			},
			ExpectedRes:           []model.Village(nil),
			MockSelectVillageErr:  errors.New("error"),
			MockSelectVillageRes:  sqlmock.NewRows([]string{"id", "name", "district_id"}),
			MockSelectDistrictErr: nil,
			MockSelectDistrictRes: sqlmock.NewRows([]string{"id", "name", "regency_id"}).
				AddRow(1, "districtA", 1),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			selectVillage := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`district_id` = ?"))
			selectVillage.WillReturnRows(v.MockSelectVillageRes)
			selectVillage.WillReturnError(v.MockSelectVillageErr)

			selectDistrict := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?"))
			selectDistrict.WillReturnRows(v.MockSelectDistrictRes)
			selectDistrict.WillReturnError(v.MockSelectDistrictErr)

			var ctx context.Context
			res, err := s.repository.FindVillage(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteRegionRepository) TestCheckIsImported() {
	testCase := []struct {
		Name              string
		ExpectedErr       error
		ExpectedRes       bool
		Body              interface{}
		MockCheckModelErr error
		MockCheckModelRes *sqlmock.Rows
	}{
		// use model province for example
		{
			Name:              "not imported",
			ExpectedErr:       nil,
			ExpectedRes:       false,
			Body:              model.Province{},
			MockCheckModelErr: gorm.ErrRecordNotFound,
			MockCheckModelRes: sqlmock.NewRows([]string{"id", "name"}),
		},
		{
			Name:              "error",
			ExpectedErr:       errors.New("error"),
			ExpectedRes:       true,
			Body:              model.Province{},
			MockCheckModelErr: errors.New("error"),
			MockCheckModelRes: sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			checkModel := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` ORDER BY `provinces`.`id` LIMIT 1"))
			checkModel.WillReturnRows(v.MockCheckModelRes)
			checkModel.WillReturnError(v.MockCheckModelErr)

			res, err := s.repository.CheckIsImported(v.Body)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteRegionRepository(t *testing.T) {
	suite.Run(t, new(suiteRegionRepository))
}
