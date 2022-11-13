package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	checkpointRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/repository/mock"
	userRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
)

type suiteCheckpointService struct {
	suite.Suite
	checkpointRepositoryMock *checkpointRepositoryMock.CheckpointRepositoryMock
	userRepositoryMock       *userRepositoryMock.UserRepositoryMock
	checkpointService        CheckpointService
}

func (s *suiteCheckpointService) SetupSuit() {
	s.checkpointRepositoryMock = new(checkpointRepositoryMock.CheckpointRepositoryMock)
	s.userRepositoryMock = new(userRepositoryMock.UserRepositoryMock)
	s.checkpointService = NewCheckpointService(s.checkpointRepositoryMock, s.userRepositoryMock)
}

func (s *suiteCheckpointService) TearDown() {
	s.checkpointRepositoryMock = nil
	s.userRepositoryMock = nil
	s.checkpointService = nil
}

func (s *suiteCheckpointService) TestCreateCheckpoint() {
	testCase := []struct {
		Name                string
		ExpectedRes         uuid.UUID
		ExpectedErr         error
		Body                dto.CheckpointRequest
		CreateCheckpointErr error
	}{
		{
			Name:        "success",
			ExpectedRes: uuid.New(),
			ExpectedErr: nil,
			Body: dto.CheckpointRequest{
				Name:        "category",
				Description: "test",
				ProvinceID:  1,
				RegencyID:   1,
				DistrictID:  1,
				VillageID:   1,
				LatLong:     "-12, 12",
			},
			CreateCheckpointErr: nil,
		},
		{
			Name:        "error foreign key",
			ExpectedRes: uuid.Nil,
			ExpectedErr: customerrors.ErrBadRequestBody,
			Body: dto.CheckpointRequest{
				Name:        "category",
				Description: "test",
				ProvinceID:  1,
				RegencyID:   1,
				DistrictID:  1,
				VillageID:   1,
				LatLong:     "-12, 12",
			},
			CreateCheckpointErr: errors.New("Cannot add or update a child row"),
		},
		{
			Name:        "other error",
			ExpectedRes: uuid.New(),
			ExpectedErr: errors.New("error"),
			Body: dto.CheckpointRequest{
				Name:        "category",
				Description: "test",
				ProvinceID:  1,
				RegencyID:   1,
				DistrictID:  1,
				VillageID:   1,
				LatLong:     "-12, 12",
			},
			CreateCheckpointErr: errors.New("error"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.checkpointRepositoryMock.On("CreateCheckpoint").Return(v.CreateCheckpointErr)

			res, err := s.checkpointService.CreateCheckpoint(v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(len(v.ExpectedRes.String()), len(res.String()))

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointService) TestFindCheckpoints() {
	checkpointId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedRes        dto.CheckpointsResponse
		ExpectedErr        error
		FindCheckpointsErr error
		FindCheckpointsRes []model.Checkpoint
	}{
		{
			Name: "success",
			ExpectedRes: dto.CheckpointsResponse{
				{ID: checkpointId, Name: "checkpoint"},
			},
			ExpectedErr:        nil,
			FindCheckpointsErr: nil,
			FindCheckpointsRes: []model.Checkpoint{
				{ID: checkpointId, Name: "checkpoint"},
			},
		},
		{
			Name:               "error",
			ExpectedRes:        dto.CheckpointsResponse(nil),
			ExpectedErr:        errors.New("error"),
			FindCheckpointsErr: errors.New("error"),
			FindCheckpointsRes: []model.Checkpoint{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.checkpointRepositoryMock.On("FindCheckpoints").Return(v.FindCheckpointsRes, v.FindCheckpointsErr)

			res, err := s.checkpointService.FindCheckpoints(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointService) TestFindCheckpointsByUser() {
	checkpointId := uuid.New()
	userId := uuid.New()

	testCase := []struct {
		Name                         string
		ExpectedRes                  dto.CheckpointsResponse
		ExpectedErr                  error
		UserId                       string
		FindUserByIdErr              error
		FindUserByIdRes              *model.User
		FindCheckpointsByVillageErr  error
		FindCheckpointsByVillageRes  []model.Checkpoint
		FindCheckpointsByDistrictErr error
		FindCheckpointsByDistrictRes []model.Checkpoint
		FindCheckpointsByRegencyErr  error
		FindCheckpointsByRegencyRes  []model.Checkpoint
		FindCheckpointsByProvinceErr error
		FindCheckpointsByProvinceRes []model.Checkpoint
	}{
		{
			Name: "success get by the village",
			ExpectedRes: dto.CheckpointsResponse{
				{
					ID: checkpointId, Name: "checkpoint",
				},
			},
			ExpectedErr:     nil,
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{{ID: checkpointId, Name: "checkpoint"}},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name: "success get by the district",
			ExpectedRes: dto.CheckpointsResponse{
				{
					ID: checkpointId, Name: "checkpoint",
				},
			},
			ExpectedErr:     nil,
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{{ID: checkpointId, Name: "checkpoint"}},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name: "success get by the regency",
			ExpectedRes: dto.CheckpointsResponse{
				{
					ID: checkpointId, Name: "checkpoint",
				},
			},
			ExpectedErr:     nil,
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{{ID: checkpointId, Name: "checkpoint"}},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name: "success get by the province",
			ExpectedRes: dto.CheckpointsResponse{
				{
					ID: checkpointId, Name: "checkpoint",
				},
			},
			ExpectedErr:     nil,
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{{ID: checkpointId, Name: "checkpoint"}},
		},
		{
			Name:            "checkpoint not covered",
			ExpectedRes:     dto.CheckpointsResponse(nil),
			ExpectedErr:     customerrors.ErrCheckpointNotCovered,
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name:                         "invalid user id",
			ExpectedRes:                  dto.CheckpointsResponse(nil),
			ExpectedErr:                  customerrors.ErrInvalidId,
			UserId:                       "123",
			FindUserByIdErr:              nil,
			FindUserByIdRes:              &model.User{},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name:            "error find by village",
			ExpectedRes:     dto.CheckpointsResponse(nil),
			ExpectedErr:     errors.New("error find by village"),
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  errors.New("error find by village"),
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name:            "error find by district",
			ExpectedRes:     dto.CheckpointsResponse(nil),
			ExpectedErr:     errors.New("error find by district"),
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: errors.New("error find by district"),
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name:            "error find by regency",
			ExpectedRes:     dto.CheckpointsResponse(nil),
			ExpectedErr:     errors.New("error find by regency"),
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  errors.New("error find by regency"),
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: nil,
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
		{
			Name:            "error find by province",
			ExpectedRes:     dto.CheckpointsResponse(nil),
			ExpectedErr:     errors.New("error find by province"),
			UserId:          userId.String(),
			FindUserByIdErr: nil,
			FindUserByIdRes: &model.User{
				ID:   userId,
				Name: "user",
			},
			FindCheckpointsByVillageErr:  nil,
			FindCheckpointsByVillageRes:  []model.Checkpoint{},
			FindCheckpointsByDistrictErr: nil,
			FindCheckpointsByDistrictRes: []model.Checkpoint{},
			FindCheckpointsByRegencyErr:  nil,
			FindCheckpointsByRegencyRes:  []model.Checkpoint{},
			FindCheckpointsByProvinceErr: errors.New("error find by province"),
			FindCheckpointsByProvinceRes: []model.Checkpoint{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.userRepositoryMock.On("FindUserByID").Return(v.FindUserByIdRes, v.FindUserByIdErr)
			s.checkpointRepositoryMock.On("FindCheckpointByVilage").Return(v.FindCheckpointsByVillageRes, v.FindCheckpointsByVillageErr)
			s.checkpointRepositoryMock.On("FindCheckpointByDistrict").Return(v.FindCheckpointsByDistrictRes, v.FindCheckpointsByDistrictErr)
			s.checkpointRepositoryMock.On("FindCheckpointByRegency").Return(v.FindCheckpointsByRegencyRes, v.FindCheckpointsByRegencyErr)
			s.checkpointRepositoryMock.On("FindCheckpointByProvince").Return(v.FindCheckpointsByProvinceRes, v.FindCheckpointsByProvinceErr)

			res, err := s.checkpointService.FindCheckpointsByUser(v.UserId, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteCheckpointService(t *testing.T) {
	suite.Run(t, new(suiteCheckpointService))
}
