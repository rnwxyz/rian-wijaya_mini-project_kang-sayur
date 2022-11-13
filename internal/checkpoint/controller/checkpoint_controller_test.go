package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	csm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/service/mock"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	mm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/middleware/mock"
	vm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/validator/mock"
	"github.com/stretchr/testify/suite"
)

type suiteCheckpointController struct {
	suite.Suite
	checkpointServiceMock *csm.CheckpointServiceMock
	JWTServiceMock        *mm.MockJWTService
	validatorMock         *vm.CustomValidatorMock
	checkpointController  *checkpointController
	echoNew               *echo.Echo
}

func (s *suiteCheckpointController) SetupSuit() {
	s.checkpointServiceMock = new(csm.CheckpointServiceMock)
	s.JWTServiceMock = new(mm.MockJWTService)
	s.validatorMock = new(vm.CustomValidatorMock)
	s.checkpointController = NewCheckpointController(s.checkpointServiceMock, s.JWTServiceMock)
	s.echoNew = echo.New()
	s.echoNew.Validator = s.validatorMock
}

func (s *suiteCheckpointController) TearDown() {
	s.checkpointServiceMock = nil
	s.JWTServiceMock = nil
	s.validatorMock = nil
	s.checkpointController = nil
	s.echoNew = nil
}

func (s *suiteCheckpointController) TestCreateCheckpoint() {
	checkpontId := uuid.New()

	testCase := []struct {
		Name                string
		ExpectedStatus      int
		ExpectedResult      map[string]interface{}
		Body                map[string]interface{}
		JWTReturn           jwt.MapClaims
		ValidatorErr        error
		CreateCheckpointErr error
		CreateCheckpointRes uuid.UUID
	}{
		{
			Name:           "success create",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"id":      checkpontId.String(),
				"message": "new checkpoint success created",
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"description": "test",
				"province_id": 1,
				"regency_id":  1,
				"district_id": 1,
				"village_id":  1,
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:        nil,
			CreateCheckpointErr: nil,
			CreateCheckpointRes: checkpontId,
		},
		{
			Name:           "forbidden",
			ExpectedStatus: 403,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrPermission.Error(),
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"description": "test",
				"province_id": 1,
				"regency_id":  1,
				"district_id": 1,
				"village_id":  1,
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(100),
			},
			ValidatorErr:        nil,
			CreateCheckpointErr: nil,
			CreateCheckpointRes: checkpontId,
		},
		{
			Name:           "invalid body type",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"description": "test",
				"province_id": "aaa",
				"regency_id":  "aaa",
				"district_id": "aaa",
				"village_id":  "aaa",
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:        nil,
			CreateCheckpointErr: nil,
			CreateCheckpointRes: checkpontId,
		},
		{
			Name:           "bad body request",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"description": "test",
				"province_id": 1,
				"regency_id":  1,
				"district_id": 1,
				"village_id":  1,
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:        nil,
			CreateCheckpointErr: customerrors.ErrBadRequestBody,
			CreateCheckpointRes: checkpontId,
		},
		{
			Name:           "validator error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("validator error").Error(),
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"province_id": 1,
				"regency_id":  1,
				"district_id": 1,
				"village_id":  1,
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:        errors.New("validator error"),
			CreateCheckpointErr: nil,
			CreateCheckpointRes: checkpontId,
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal errror").Error(),
			},
			Body: map[string]interface{}{
				"name":        "checkpoint",
				"description": "test",
				"province_id": 1,
				"regency_id":  1,
				"district_id": 1,
				"village_id":  1,
				"lat_long":    "-12, 12",
			},
			JWTReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValidatorErr:        nil,
			CreateCheckpointErr: errors.New("internal errror"),
			CreateCheckpointRes: checkpontId,
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/checkpoints")

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.checkpointServiceMock.On("CreateCheckpoint").Return(v.CreateCheckpointRes, v.CreateCheckpointErr)
			s.validatorMock.On("Validate").Return(v.ValidatorErr)

			err = s.checkpointController.CreateCheckpoint(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}

func (s *suiteCheckpointController) TestGetCheckpoint() {
	checkpointId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedStatus     int
		ExpectedResult     map[string]interface{}
		FindCheckpointsErr error
		FindCheckpointsRes dto.CheckpointsResponse
	}{
		{
			Name:           "find checkpoint success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"description":   "",
						"district_name": "district",
						"id":            checkpointId.String(),
						"lat_long":      "",
						"name":          "checkpoint",
						"province_name": "province",
						"regency_name":  "regency",
						"village_name":  "village",
					},
				},
				"message": "get checkpoint success",
			},
			FindCheckpointsErr: nil,
			FindCheckpointsRes: dto.CheckpointsResponse{
				{
					ID:           checkpointId,
					Name:         "checkpoint",
					ProvinceName: "province",
					RegencyName:  "regency",
					DistrictName: "district",
					VillageName:  "village",
				},
			},
		},
		{
			Name:           "error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error").Error(),
			},
			FindCheckpointsErr: errors.New("error"),
			FindCheckpointsRes: dto.CheckpointsResponse{},
		},
		{
			Name:           "nil checkpoint",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data":    []interface{}{},
				"message": "get checkpoint success",
			},
			FindCheckpointsErr: nil,
			FindCheckpointsRes: dto.CheckpointsResponse{},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/checkpoints")

			// define mock
			s.checkpointServiceMock.On("FindCheckpoints").Return(v.FindCheckpointsRes, v.FindCheckpointsErr)

			err := s.checkpointController.GetCheckpoints(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}
func (s *suiteCheckpointController) TestGetCheckpointByUser() {
	checkpointId := uuid.New()
	userId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedStatus     int
		ExpectedResult     map[string]interface{}
		JwtReturn          jwt.MapClaims
		FindCheckpointsErr error
		FindCheckpointsRes dto.CheckpointsResponse
	}{
		{
			Name:           "find checkpoint success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"description":   "",
						"district_name": "district",
						"id":            checkpointId.String(),
						"lat_long":      "",
						"name":          "checkpoint",
						"province_name": "province",
						"regency_name":  "regency",
						"village_name":  "village",
					},
				},
				"message": "get checkpoint success",
			},
			JwtReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindCheckpointsErr: nil,
			FindCheckpointsRes: dto.CheckpointsResponse{
				{
					ID:           checkpointId,
					Name:         "checkpoint",
					ProvinceName: "province",
					RegencyName:  "regency",
					DistrictName: "district",
					VillageName:  "village",
				},
			},
		},
		{
			Name:           "checkpoint not covered",
			ExpectedStatus: 404,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrCheckpointNotCovered.Error(),
			},
			JwtReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindCheckpointsErr: customerrors.ErrCheckpointNotCovered,
			FindCheckpointsRes: dto.CheckpointsResponse{},
		},
		{
			Name:           "error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error").Error(),
			},
			JwtReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindCheckpointsErr: errors.New("error"),
			FindCheckpointsRes: dto.CheckpointsResponse{},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/checkpoints")

			// define mock
			s.checkpointServiceMock.On("FindCheckpointsByUser").Return(v.FindCheckpointsRes, v.FindCheckpointsErr)
			s.JWTServiceMock.On("GetClaims").Return(v.JwtReturn)

			err := s.checkpointController.GetCheckpointByUser(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			s.TearDown()
		})
	}
}

func TestCheckpointController(t *testing.T) {
	suite.Run(t, new(suiteCheckpointController))
}
