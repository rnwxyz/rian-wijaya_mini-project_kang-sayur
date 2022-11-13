package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/service"
	rsm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/service/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
)

type suiteRegionController struct {
	suite.Suite
	regionServiceMock *rsm.RegionServiceMock
	regionController  *regionController
}

func newregionControllerMock(service service.RegionService) *regionController {
	return &regionController{
		service: service,
	}
}

func (s *suiteRegionController) SetupSuit() {
	s.regionServiceMock = new(rsm.RegionServiceMock)
	s.regionController = newregionControllerMock(s.regionServiceMock)
}

func (s *suiteRegionController) TearDown() {
	s.regionServiceMock = nil
	s.regionController = nil
}

func (s *suiteRegionController) TestGetProvince() {
	testCase := []struct {
		Name            string
		ExpectedStatus  int
		FindProvinceErr error
		FindProvinceRes []model.Province
		ExpectedResult  map[string]interface{}
	}{
		{
			Name:            "success",
			ExpectedStatus:  200,
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
			ExpectedResult: map[string]interface{}{
				"message": "get provinces success",
				"data": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "provinceA",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "provinceB",
					},
				},
			},
		},
		{
			Name:            "internal server error",
			ExpectedStatus:  500,
			FindProvinceErr: errors.New("internal server error"),
			FindProvinceRes: []model.Province{},
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal server error").Error(),
			},
		},
		{
			Name:            "nil province",
			ExpectedStatus:  200,
			FindProvinceErr: nil,
			FindProvinceRes: []model.Province(nil),
			ExpectedResult: map[string]interface{}{
				"message": "get provinces success",
				"data":    nil,
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/region/province")

			// define mock
			s.regionServiceMock.On("FindProvince").Return(v.FindProvinceRes, v.FindProvinceErr)

			err := s.regionController.GetProvince(ctx)
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
func (s *suiteRegionController) TestGetRegency() {
	testCase := []struct {
		Name           string
		ProvinceId     string
		ExpectedStatus int
		FindRegencyErr error
		FindRegencyRes dto.RegenciesResponse
		ExpectedResult map[string]interface{}
	}{
		{
			Name:           "success",
			ProvinceId:     "1",
			ExpectedStatus: 200,
			FindRegencyErr: nil,
			FindRegencyRes: dto.RegenciesResponse{
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
			ExpectedResult: map[string]interface{}{
				"message": "get regencies success",
				"data": []interface{}{
					map[string]interface{}{
						"id":            float64(1),
						"name":          "regencyA",
						"province_id":   float64(1),
						"province_name": "provinceA",
					},
					map[string]interface{}{
						"id":            float64(2),
						"name":          "regencyB",
						"province_id":   float64(1),
						"province_name": "provinceA",
					},
				},
			},
		},
		{
			Name:           "internal server error",
			ProvinceId:     "1",
			ExpectedStatus: 500,
			FindRegencyErr: errors.New("internal server error"),
			FindRegencyRes: dto.RegenciesResponse{},
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal server error").Error(),
			},
		},
		{
			Name:           "invalid id param",
			ProvinceId:     "1",
			ExpectedStatus: 400,
			FindRegencyErr: customerrors.ErrInvalidId,
			FindRegencyRes: dto.RegenciesResponse{},
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrInvalidId.Error(),
			},
		},
		{
			Name:           "nil regency",
			ProvinceId:     "1",
			ExpectedStatus: 200,
			FindRegencyErr: nil,
			FindRegencyRes: dto.RegenciesResponse(nil),
			ExpectedResult: map[string]interface{}{
				"message": "get regencies success",
				"data":    nil,
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/region/Regency/:province_id")
			ctx.SetParamNames("province_id")
			ctx.SetParamValues(v.ProvinceId)

			// define mock
			s.regionServiceMock.On("FindRegency").Return(v.FindRegencyRes, v.FindRegencyErr)

			err := s.regionController.GetRegency(ctx)
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
func (s *suiteRegionController) TestGetDistrict() {
	testCase := []struct {
		Name            string
		RegencyId       string
		ExpectedStatus  int
		FindDistrictErr error
		FindDistrictRes dto.DistrictsResponse
		ExpectedResult  map[string]interface{}
	}{
		{
			Name:            "success",
			RegencyId:       "1",
			ExpectedStatus:  200,
			FindDistrictErr: nil,
			FindDistrictRes: dto.DistrictsResponse{
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
			ExpectedResult: map[string]interface{}{
				"message": "get districts success",
				"data": []interface{}{
					map[string]interface{}{
						"id":           float64(1),
						"name":         "districtA",
						"regency_id":   float64(1),
						"regency_name": "regencyA",
					},
					map[string]interface{}{
						"id":           float64(2),
						"name":         "districtB",
						"regency_id":   float64(1),
						"regency_name": "regencyA",
					},
				},
			},
		},
		{
			Name:            "internal server error",
			RegencyId:       "1",
			ExpectedStatus:  500,
			FindDistrictErr: errors.New("internal server error"),
			FindDistrictRes: dto.DistrictsResponse{},
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal server error").Error(),
			},
		},
		{
			Name:            "invalid id param",
			RegencyId:       "1",
			ExpectedStatus:  400,
			FindDistrictErr: customerrors.ErrInvalidId,
			FindDistrictRes: dto.DistrictsResponse{},
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrInvalidId.Error(),
			},
		},
		{
			Name:            "nil district",
			RegencyId:       "1",
			ExpectedStatus:  200,
			FindDistrictErr: nil,
			FindDistrictRes: dto.DistrictsResponse(nil),
			ExpectedResult: map[string]interface{}{
				"message": "get districts success",
				"data":    nil,
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/region/district/:regency_id")
			ctx.SetParamNames("regency_id")
			ctx.SetParamValues(v.RegencyId)

			// define mock
			s.regionServiceMock.On("FindDistrict").Return(v.FindDistrictRes, v.FindDistrictErr)

			err := s.regionController.GetDistrict(ctx)
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
func (s *suiteRegionController) TestGetVillage() {
	testCase := []struct {
		Name           string
		DistrictId     string
		ExpectedStatus int
		FindVillageErr error
		FindVillageRes dto.VillagesResponse
		ExpectedResult map[string]interface{}
	}{
		{
			Name:           "success",
			DistrictId:     "1",
			ExpectedStatus: 200,
			FindVillageErr: nil,
			FindVillageRes: dto.VillagesResponse{
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
			ExpectedResult: map[string]interface{}{
				"message": "get villages success",
				"data": []interface{}{
					map[string]interface{}{
						"id":            float64(1),
						"name":          "villageA",
						"district_id":   float64(1),
						"district_name": "districtA",
					},
					map[string]interface{}{
						"id":            float64(2),
						"name":          "villageB",
						"district_id":   float64(1),
						"district_name": "districtA",
					},
				},
			},
		},
		{
			Name:           "internal server error",
			DistrictId:     "1",
			ExpectedStatus: 500,
			FindVillageErr: errors.New("internal server error"),
			FindVillageRes: dto.VillagesResponse{},
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal server error").Error(),
			},
		},
		{
			Name:           "invalid id param",
			DistrictId:     "1",
			ExpectedStatus: 400,
			FindVillageErr: customerrors.ErrInvalidId,
			FindVillageRes: dto.VillagesResponse{},
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrInvalidId.Error(),
			},
		},
		{
			Name:           "nil village",
			DistrictId:     "1",
			ExpectedStatus: 200,
			FindVillageErr: nil,
			FindVillageRes: dto.VillagesResponse(nil),
			ExpectedResult: map[string]interface{}{
				"message": "get villages success",
				"data":    nil,
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := echo.New()
			ctx := c.NewContext(r, w)
			ctx.SetPath("/region/village/:district_id")
			ctx.SetParamNames("district_id")
			ctx.SetParamValues(v.DistrictId)

			// define mock
			s.regionServiceMock.On("FindVillage").Return(v.FindVillageRes, v.FindVillageErr)

			err := s.regionController.GetVillage(ctx)
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

func (s *suiteRegionController) TestInitRoute() {
	group := echo.New().Group("/api/v1")
	s.NotPanics(func() {
		s.regionController.InitRoute(group)
	})
}

func TestRegionController(t *testing.T) {
	suite.Run(t, new(suiteRegionController))
}
