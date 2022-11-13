package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteCheckpointRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *checkpointRepositoryImpl
}

func (s *suiteCheckpointRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &checkpointRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteCheckpointRepository) TearDown() {
	s.mock = nil
	s.repository = nil
}

func (s *suiteCheckpointRepository) TestCreateCheckpoint() {
	checkpointId := uuid.New()

	testCase := []struct {
		Name        string
		Body        model.Checkpoint
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success create",
			Body: model.Checkpoint{
				ID:   checkpointId,
				Name: "checkpoint",
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "error",
			Body: model.Checkpoint{
				ID:   checkpointId,
				Name: "checkpoint",
			},
			ExpectedErr: errors.New("error"),
			MockReturn:  errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `checkpoints` (`id`,`created_at`,`updated_at`,`deleted_at`,`name`,`description`,`province_id`,`regency_id`,`district_id`,`village_id`,`lat_long`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}

			err := s.repository.CreateCheckpoint(&v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointRepository) TestFindCheckpoints() {
	checkpointId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedErr        error
		ExpectedRes        []model.Checkpoint
		FindCheckpointErr  error
		FindCheckpointRes  *sqlmock.Rows
		PreloadProvinceRes *sqlmock.Rows
		PreloadRegencyRes  *sqlmock.Rows
		PreloadDistrictRes *sqlmock.Rows
		PreloadVillageRes  *sqlmock.Rows
	}{
		{
			Name:        "success find checkpoints",
			ExpectedErr: nil,
			ExpectedRes: []model.Checkpoint{
				{
					ID:         checkpointId,
					Name:       "checkpoint",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "province",
					},
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regency",
					},
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "district",
					},
					VillageID: 1,
					Village: model.Village{
						ID:   1,
						Name: "village",
					},
					LatLong: "-12, 12",
				},
			},
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}).AddRow(checkpointId, "checkpoint", 1, 1, 1, 1, "-12, 12"),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "province"),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "regency"),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "district"),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "village"),
		},
		{
			Name:               "error",
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  errors.New("error"),
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE `checkpoints`.`deleted_at` IS NULL")).WillReturnError(v.ExpectedErr).WillReturnRows(v.FindCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?")).WillReturnRows(v.PreloadDistrictRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?")).WillReturnRows(v.PreloadProvinceRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?")).WillReturnRows(v.PreloadRegencyRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`id` = ?")).WillReturnRows(v.PreloadVillageRes)

			res, err := s.repository.FindCheckpoints(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointRepository) TestFindCheckpointByVillage() {
	checkpointId := uuid.New()
	villageId := uint(1)

	testCase := []struct {
		Name               string
		ExpectedErr        error
		User               model.User
		ExpectedRes        []model.Checkpoint
		FindCheckpointErr  error
		FindCheckpointRes  *sqlmock.Rows
		PreloadProvinceRes *sqlmock.Rows
		PreloadRegencyRes  *sqlmock.Rows
		PreloadDistrictRes *sqlmock.Rows
		PreloadVillageRes  *sqlmock.Rows
	}{
		{
			Name:        "success find checkpoints",
			ExpectedErr: nil,
			User: model.User{
				VillageID: &villageId,
			},
			ExpectedRes: []model.Checkpoint{
				{
					ID:         checkpointId,
					Name:       "checkpoint",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "province",
					},
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regency",
					},
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "district",
					},
					VillageID: 1,
					Village: model.Village{
						ID:   1,
						Name: "village",
					},
					LatLong: "-12, 12",
				},
			},
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}).AddRow(checkpointId, "checkpoint", 1, 1, 1, 1, "-12, 12"),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "province"),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "regency"),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "district"),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "village"),
		},
		{
			Name: "error",
			User: model.User{
				VillageID: &villageId,
			},
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  errors.New("error"),
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
		{
			Name: "user dont have village id",
			User: model.User{
				VillageID: nil,
			},
			ExpectedErr:        nil,
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE village_id = ? AND `checkpoints`.`deleted_at` IS NULL")).WillReturnError(v.ExpectedErr).WillReturnRows(v.FindCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?")).WillReturnRows(v.PreloadDistrictRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?")).WillReturnRows(v.PreloadProvinceRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?")).WillReturnRows(v.PreloadRegencyRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`id` = ?")).WillReturnRows(v.PreloadVillageRes)

			res, err := s.repository.FindCheckpointByVilage(v.User, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointRepository) TestFindCheckpointByDistrict() {
	checkpointId := uuid.New()
	districtId := uint(1)

	testCase := []struct {
		Name               string
		ExpectedErr        error
		User               model.User
		ExpectedRes        []model.Checkpoint
		FindCheckpointErr  error
		FindCheckpointRes  *sqlmock.Rows
		PreloadProvinceRes *sqlmock.Rows
		PreloadRegencyRes  *sqlmock.Rows
		PreloadDistrictRes *sqlmock.Rows
		PreloadVillageRes  *sqlmock.Rows
	}{
		{
			Name:        "success find checkpoints",
			ExpectedErr: nil,
			User: model.User{
				DistrictID: &districtId,
			},
			ExpectedRes: []model.Checkpoint{
				{
					ID:         checkpointId,
					Name:       "checkpoint",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "province",
					},
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regency",
					},
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "district",
					},
					VillageID: 1,
					Village: model.Village{
						ID:   1,
						Name: "village",
					},
					LatLong: "-12, 12",
				},
			},
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}).AddRow(checkpointId, "checkpoint", 1, 1, 1, 1, "-12, 12"),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "province"),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "regency"),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "district"),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "village"),
		},
		{
			Name: "error",
			User: model.User{
				DistrictID: &districtId,
			},
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  errors.New("error"),
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
		{
			Name: "user dont have village id",
			User: model.User{
				DistrictID: nil,
			},
			ExpectedErr:        nil,
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE district_id = ? AND `checkpoints`.`deleted_at` IS NULL")).WillReturnError(v.ExpectedErr).WillReturnRows(v.FindCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?")).WillReturnRows(v.PreloadDistrictRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?")).WillReturnRows(v.PreloadProvinceRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?")).WillReturnRows(v.PreloadRegencyRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`id` = ?")).WillReturnRows(v.PreloadVillageRes)

			res, err := s.repository.FindCheckpointByDistrict(v.User, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointRepository) TestFindCheckpointByRegency() {
	checkpointId := uuid.New()
	regencyId := uint(1)

	testCase := []struct {
		Name               string
		ExpectedErr        error
		User               model.User
		ExpectedRes        []model.Checkpoint
		FindCheckpointErr  error
		FindCheckpointRes  *sqlmock.Rows
		PreloadProvinceRes *sqlmock.Rows
		PreloadRegencyRes  *sqlmock.Rows
		PreloadDistrictRes *sqlmock.Rows
		PreloadVillageRes  *sqlmock.Rows
	}{
		{
			Name:        "success find checkpoints",
			ExpectedErr: nil,
			User: model.User{
				RegencyID: &regencyId,
			},
			ExpectedRes: []model.Checkpoint{
				{
					ID:         checkpointId,
					Name:       "checkpoint",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "province",
					},
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regency",
					},
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "district",
					},
					VillageID: 1,
					Village: model.Village{
						ID:   1,
						Name: "village",
					},
					LatLong: "-12, 12",
				},
			},
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}).AddRow(checkpointId, "checkpoint", 1, 1, 1, 1, "-12, 12"),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "province"),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "regency"),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "district"),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "village"),
		},
		{
			Name: "error",
			User: model.User{
				RegencyID: &regencyId,
			},
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  errors.New("error"),
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
		{
			Name: "user dont have village id",
			User: model.User{
				RegencyID: nil,
			},
			ExpectedErr:        nil,
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE regency_id = ? AND `checkpoints`.`deleted_at` IS NULL")).WillReturnError(v.ExpectedErr).WillReturnRows(v.FindCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?")).WillReturnRows(v.PreloadDistrictRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?")).WillReturnRows(v.PreloadProvinceRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?")).WillReturnRows(v.PreloadRegencyRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`id` = ?")).WillReturnRows(v.PreloadVillageRes)

			res, err := s.repository.FindCheckpointByRegency(v.User, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointRepository) TestFindCheckpointByProvince() {
	checkpointId := uuid.New()
	provinceId := uint(1)

	testCase := []struct {
		Name               string
		ExpectedErr        error
		User               model.User
		ExpectedRes        []model.Checkpoint
		FindCheckpointErr  error
		FindCheckpointRes  *sqlmock.Rows
		PreloadProvinceRes *sqlmock.Rows
		PreloadRegencyRes  *sqlmock.Rows
		PreloadDistrictRes *sqlmock.Rows
		PreloadVillageRes  *sqlmock.Rows
	}{
		{
			Name:        "success find checkpoints",
			ExpectedErr: nil,
			User: model.User{
				ProvinceID: &provinceId,
			},
			ExpectedRes: []model.Checkpoint{
				{
					ID:         checkpointId,
					Name:       "checkpoint",
					ProvinceID: 1,
					Province: model.Province{
						ID:   1,
						Name: "province",
					},
					RegencyID: 1,
					Regency: model.Regency{
						ID:   1,
						Name: "regency",
					},
					DistrictID: 1,
					District: model.District{
						ID:   1,
						Name: "district",
					},
					VillageID: 1,
					Village: model.Village{
						ID:   1,
						Name: "village",
					},
					LatLong: "-12, 12",
				},
			},
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}).AddRow(checkpointId, "checkpoint", 1, 1, 1, 1, "-12, 12"),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "province"),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "regency"),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "district"),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "village"),
		},
		{
			Name: "error",
			User: model.User{
				ProvinceID: &provinceId,
			},
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  errors.New("error"),
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
		{
			Name: "user dont have village id",
			User: model.User{
				ProvinceID: nil,
			},
			ExpectedErr:        nil,
			ExpectedRes:        []model.Checkpoint(nil),
			FindCheckpointErr:  nil,
			FindCheckpointRes:  sqlmock.NewRows([]string{"id", "name", "province_id", "regency_id", "district_id", "village_id", "lat_long"}),
			PreloadProvinceRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadRegencyRes:  sqlmock.NewRows([]string{"id", "name"}),
			PreloadDistrictRes: sqlmock.NewRows([]string{"id", "name"}),
			PreloadVillageRes:  sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE province_id = ? AND `checkpoints`.`deleted_at` IS NULL")).WillReturnError(v.ExpectedErr).WillReturnRows(v.FindCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `districts` WHERE `districts`.`id` = ?")).WillReturnRows(v.PreloadDistrictRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provinces` WHERE `provinces`.`id` = ?")).WillReturnRows(v.PreloadProvinceRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `regencies` WHERE `regencies`.`id` = ?")).WillReturnRows(v.PreloadRegencyRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `villages` WHERE `villages`.`id` = ?")).WillReturnRows(v.PreloadVillageRes)

			res, err := s.repository.FindCheckpointByProvince(v.User, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteCheckpointRepository(t *testing.T) {
	suite.Run(t, new(suiteCheckpointRepository))
}
