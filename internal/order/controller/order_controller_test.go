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
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	osm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	mm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/middleware/mock"
	qrm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/qrcode/mock"
	vm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/validator/mock"
	"github.com/stretchr/testify/suite"
)

type suiteOrderController struct {
	suite.Suite
	orderServiceMock *osm.OrderServiceMock
	JWTServiceMock   *mm.MockJWTService
	QrCodeMock       *qrm.QRCodeMock
	orderController  *orderController
	validatorMock    *vm.CustomValidatorMock
	echoNew          *echo.Echo
}

func newOrderControllerMock(service service.OrderService, jwt JWTService, qr QRCode) *orderController {
	return &orderController{
		service:    service,
		jwtService: jwt,
		qrCode:     qr,
	}
}

func (s *suiteOrderController) SetupSuit() {
	s.orderServiceMock = new(osm.OrderServiceMock)
	s.JWTServiceMock = new(mm.MockJWTService)
	s.QrCodeMock = new(qrm.QRCodeMock)
	s.validatorMock = new(vm.CustomValidatorMock)
	s.orderController = newOrderControllerMock(s.orderServiceMock, s.JWTServiceMock, s.QrCodeMock)
	s.echoNew = echo.New()
	s.echoNew.Validator = s.validatorMock
}

func (s *suiteOrderController) TearDown() {
	s.orderServiceMock = nil
	s.JWTServiceMock = nil
	s.QrCodeMock = nil
	s.orderController = nil
	s.validatorMock = nil
	s.echoNew = nil
}

func (s *suiteOrderController) TestCreateOrder() {
	userId := uuid.New()
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		Body           map[string]interface{}
		JWTReturn      jwt.MapClaims
		ValidatorErr   error
		CreateOrderErr error
		CreateOrderRes *dto.NewOrder
	}{
		{
			Name:           "success create",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": map[string]interface{}{
					"order_id":     orderId.String(),
					"redirect_url": "https://test/miodtrans.com",
				},
				"message": "new create order success created",
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: nil,
			CreateOrderRes: &dto.NewOrder{
				OrderID:     orderId,
				RedirectURL: "https://test/miodtrans.com",
			},
		},
		{
			Name:           "error invalid body type",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": 123,
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: nil,
			CreateOrderRes: &dto.NewOrder{},
		},
		{
			Name:           "validation error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("error validation").Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   errors.New("error validation"),
			CreateOrderErr: nil,
			CreateOrderRes: &dto.NewOrder{},
		},
		{
			Name:           "create order invalid id error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrInvalidId.Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: customerrors.ErrInvalidId,
			CreateOrderRes: &dto.NewOrder{},
		},
		{
			Name:           "create order error qty exceeds stock or less that 1",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrQtyOrder.Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: customerrors.ErrQtyOrder,
			CreateOrderRes: &dto.NewOrder{},
		},
		{
			Name:           "create order bad body request",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrBadRequestBody.Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: customerrors.ErrBadRequestBody,
			CreateOrderRes: &dto.NewOrder{},
		},
		{
			Name:           "create order internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal error").Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"order": []interface{}{
					map[string]interface{}{
						"item_id": 1,
						"qty":     10,
					},
				},
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			ValidatorErr:   nil,
			CreateOrderErr: errors.New("internal error"),
			CreateOrderRes: &dto.NewOrder{},
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
			ctx.SetPath("/orders")

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.validatorMock.On("Validate").Return(v.ValidatorErr)
			s.orderServiceMock.On("CreateOrder").Return(v.CreateOrderRes, v.CreateOrderErr)

			err = s.orderController.CreateOrder(ctx)
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
func (s *suiteOrderController) TestGetOrder() {
	userId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		JWTReturn      jwt.MapClaims
		FindOrderErr   error
		FindOrderRes   dto.OrdersResponse
	}{
		{
			Name:           "success get order",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"checkpoint_id": "00000000-0000-0000-0000-000000000000", "checkpoint_name": "",
						"created_at":    "0001-01-01T00:00:00Z",
						"expired_order": "0001-01-01T00:00:00Z",
						"grand_total":   float64(0),
						"id":            orderId.String(),
						"shipping_cost": float64(0),
						"status_order":  "", "total_price": float64(0),
						"user_id":   userId.String(),
						"user_name": "",
					},
				},
				"message": "get orders success",
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderErr: nil,
			FindOrderRes: dto.OrdersResponse{
				{
					ID:     orderId,
					UserID: userId,
				},
			},
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("inetrnal error").Error(),
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderErr: errors.New("inetrnal error"),
			FindOrderRes: dto.OrdersResponse{},
		},
		{
			Name:           "nil or empty order",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data":    []interface{}{},
				"message": "get orders success",
			},
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderErr: nil,
			FindOrderRes: dto.OrdersResponse{},
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
			ctx.SetPath("/orders")

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.orderServiceMock.On("FindOrder").Return(v.FindOrderRes, v.FindOrderErr)

			err := s.orderController.GetOrder(ctx)
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
func (s *suiteOrderController) TestGetOrderDetail() {
	userId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name               string
		ExpectedStatus     int
		ExpectedResult     map[string]interface{}
		ParamId            string
		JWTReturn          jwt.MapClaims
		FindOrderDetailErr error
		FindOrderDetailRes *dto.OrderWithDetailResponse
	}{
		{
			Name:           "success get detail",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": map[string]interface{}{
					"checkpoint_id":   "00000000-0000-0000-0000-000000000000",
					"checkpoint_name": "",
					"code":            "",
					"created_at":      "0001-01-01T00:00:00Z",
					"expired_order":   "0001-01-01T00:00:00Z",
					"grand_total":     float64(0),
					"id":              orderId.String(),
					"order_detail":    interface{}(nil),
					"shipping_cost":   float64(0),
					"status_order":    "",
					"total_price":     float64(0),
				},
				"message": "get orders success",
			},
			ParamId: orderId.String(),
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderDetailErr: nil,
			FindOrderDetailRes: &dto.OrderWithDetailResponse{
				ID: orderId,
			},
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal error").Error(),
			},
			ParamId: orderId.String(),
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderDetailErr: errors.New("internal error"),
			FindOrderDetailRes: &dto.OrderWithDetailResponse{},
		},
		{
			Name:           "nil or empty order",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"data": map[string]interface{}{
					"checkpoint_id":   "00000000-0000-0000-0000-000000000000",
					"checkpoint_name": "",
					"code":            "",
					"created_at":      "0001-01-01T00:00:00Z",
					"expired_order":   "0001-01-01T00:00:00Z",
					"grand_total":     float64(0),
					"id":              "00000000-0000-0000-0000-000000000000",
					"order_detail":    interface{}(nil),
					"shipping_cost":   float64(0),
					"status_order":    "",
					"total_price":     float64(0),
				},
				"message": "get orders success",
			},
			ParamId: orderId.String(),
			JWTReturn: jwt.MapClaims{
				"user_id": userId.String(),
			},
			FindOrderDetailErr: nil,
			FindOrderDetailRes: &dto.OrderWithDetailResponse{},
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
			ctx.SetPath("/orders/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(orderId.String())

			// define mock
			s.JWTServiceMock.On("GetClaims").Return(v.JWTReturn)
			s.orderServiceMock.On("FindOrderDetail").Return(v.FindOrderDetailRes, v.FindOrderDetailErr)

			err := s.orderController.GetOrderDetail(ctx)
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
func (s *suiteOrderController) TestCelcelOrder() {
	orderId := uuid.New()

	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		ParamId        string
		CencelOrderErr error
	}{
		{
			Name:           "success cencel order",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"message": "order cencelled",
			},
			ParamId:        orderId.String(),
			CencelOrderErr: nil,
		},
		{
			Name:           "internal server error",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("internal error").Error(),
			},
			ParamId:        orderId.String(),
			CencelOrderErr: errors.New("internal error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPut, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/orders/cencel/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamId)

			// define mock
			s.orderServiceMock.On("CencelOder").Return(v.CencelOrderErr)

			err := s.orderController.CencelOrder(ctx)
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
func (s *suiteOrderController) TestTakeOrder() {
	checkpointId := uuid.New()

	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		Body           map[string]interface{}
		JwtReturn      jwt.MapClaims
		ValiodatorErr  error
		TakeOrderErr   error
	}{
		{
			Name:           "success take order",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"message": "success take order",
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"code":          "qwert",
			},
			JwtReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValiodatorErr: nil,
			TakeOrderErr:  nil,
		},
		{
			Name:           "forbidden",
			ExpectedStatus: 403,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrPermission.Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"code":          "qwert",
			},
			JwtReturn: jwt.MapClaims{
				"role_id": float64(2),
			},
			ValiodatorErr: nil,
			TakeOrderErr:  nil,
		},
		{
			Name:           "validator error",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("validator error").Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"code":          "qwert",
			},
			JwtReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValiodatorErr: errors.New("validator error"),
			TakeOrderErr:  nil,
		},
		{
			Name:           "error when take order",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("take order error").Error(),
			},
			Body: map[string]interface{}{
				"checkpoint_id": checkpointId.String(),
				"code":          "qwert",
			},
			JwtReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			ValiodatorErr: nil,
			TakeOrderErr:  errors.New("take order error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/orders/takeorder")

			// define mock
			s.orderServiceMock.On("TakeOrder").Return(v.TakeOrderErr)
			s.validatorMock.On("Validate").Return(v.ValiodatorErr)
			s.JWTServiceMock.On("GetClaims").Return(v.JwtReturn)

			err := s.orderController.TakeOrder(ctx)
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
func (s *suiteOrderController) TestOrderReady() {
	orderId := uuid.New()

	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		OrderId        string
		JwtReturn      jwt.MapClaims
		OrderReadyErr  error
	}{
		{
			Name:           "order ready success",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"message": "order ready in checkpoint",
			},
			OrderId: orderId.String(),
			JwtReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			OrderReadyErr: nil,
		},
		{
			Name:           "forbidden",
			ExpectedStatus: 403,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrPermission.Error(),
			},
			OrderId: orderId.String(),
			JwtReturn: jwt.MapClaims{
				"role_id": float64(2),
			},
			OrderReadyErr: nil,
		},
		{
			Name:           "error when take order",
			ExpectedStatus: 500,
			ExpectedResult: map[string]interface{}{
				"message": errors.New("take order error").Error(),
			},
			OrderId: orderId.String(),
			JwtReturn: jwt.MapClaims{
				"role_id": float64(constants.Role_admin),
			},
			OrderReadyErr: errors.New("take order error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c := s.echoNew
			ctx := c.NewContext(r, w)
			ctx.SetPath("/orders/ready/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.OrderId)

			// define mock
			s.orderServiceMock.On("OrderReady").Return(v.OrderReadyErr)
			s.JWTServiceMock.On("GetClaims").Return(v.JwtReturn)

			err := s.orderController.OrderReady(ctx)
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
func (s *suiteOrderController) TestGetQRCode() {
	testCase := []struct {
		Name           string
		ExpectedStatus int
		ExpectedResult map[string]interface{}
		HashCode       string
		GenerateQrErr  error
		GenerateQrRes  []byte
	}{
		{
			Name:           "success get qrcode",
			ExpectedStatus: 200,
			ExpectedResult: map[string]interface{}{
				"message": "get qrcode success",
				"data": map[string]interface{}{
					"base64": "AQID",
				},
			},
			HashCode:      "abc",
			GenerateQrErr: nil,
			GenerateQrRes: []byte{
				1, 2, 3,
			},
		},
		{
			Name:           "error get qrcode",
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": customerrors.ErrGenerateQR.Error(),
			},
			HashCode:      "abc",
			GenerateQrErr: errors.New("error"),
			GenerateQrRes: []byte{},
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
			ctx.SetPath("/orders/qr/:hash_code")
			ctx.SetParamNames("hash_code")
			ctx.SetParamValues(v.HashCode)

			// define mock
			s.QrCodeMock.On("GenerateQRCode").Return(v.GenerateQrRes, v.GenerateQrErr)

			err := s.orderController.GetQRCode(ctx)
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

func TestOrderController(t *testing.T) {
	suite.Run(t, new(suiteOrderController))
}
