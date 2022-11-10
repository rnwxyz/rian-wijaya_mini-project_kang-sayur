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
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	usm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	utm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/mock"
	"github.com/stretchr/testify/suite"
)

type suiteUserController struct {
	suite.Suite
	userServiceMock *usm.UserServiceMock
	JWTServiceMock  *utm.MockJWTService
	validaorMock    *utm.CustomValidatorMock
	userController  *userController
	echoNew         *echo.Echo
}

func newUserControllerMock(service service.UserService, jwt JWTService) *userController {
	return &userController{
		service:    service,
		jwtService: jwt,
	}
}

func (s *suiteUserController) SetupTest() {
	s.userServiceMock = new(usm.UserServiceMock)
	s.JWTServiceMock = new(utm.MockJWTService)
	s.validaorMock = new(utm.CustomValidatorMock)
	s.userController = newUserControllerMock(s.userServiceMock, s.JWTServiceMock)
	s.echoNew = echo.New()
	s.echoNew.Validator = s.validaorMock
}

func (s *suiteUserController) TestSignUp() {
	// uuid statis
	varUUID := uuid.New()

	testCase := []struct {
		Name           string
		Body           map[string]interface{}
		ValidatorErr   error
		ExpectedStatus int
		CreateUserRes  uuid.UUID
		CreateUserErr  error
		ExpectedResult map[string]interface{}
	}{
		{
			Name: "success",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 200,
			CreateUserRes:  varUUID,
			CreateUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"id":      varUUID.String(),
				"message": "sign up success",
			},
		},
		{
			Name: "body type invalid",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": 123,
			},
			ValidatorErr:   nil,
			ExpectedStatus: 400,
			CreateUserRes:  uuid.Nil,
			CreateUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrBadRequestBody.Error(),
			},
		},
		{
			Name: "body invalid",
			Body: map[string]interface{}{
				"email": "test@gmail.com",
			},
			ValidatorErr:   errors.New("invalid"),
			ExpectedStatus: 400,
			CreateUserRes:  uuid.Nil,
			CreateUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"message": "invalid",
			},
		},
		{
			Name: "duplicate email",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 400,
			CreateUserRes:  uuid.Nil,
			CreateUserErr:  utils.ErrEmailAlredyExist,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrEmailAlredyExist.Error(),
			},
		},
		{
			Name: "internal server error",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 500,
			CreateUserRes:  uuid.Nil,
			CreateUserErr:  errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/signup")

			// define mock
			cus := s.userServiceMock.On("CreateUser").Return(v.CreateUserRes, v.CreateUserErr)
			val := s.validaorMock.On("Validate").Return(v.ValidatorErr)

			err = s.userController.SignUp(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			cus.Unset()
			val.Unset()
		})
	}
}

func (s *suiteUserController) TestLogin() {
	testCase := []struct {
		Name           string
		Body           map[string]interface{}
		ValidatorErr   error
		ExpectedStatus int
		LoginRes       string
		LoginErr       error
		ExpectedResult map[string]interface{}
	}{
		{
			Name: "success",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 200,
			LoginRes:       "123",
			LoginErr:       nil,
			ExpectedResult: map[string]interface{}{
				"message": "login success",
				"token":   "123",
			},
		},
		{
			Name: "body type invalid",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": 123,
			},
			ValidatorErr:   nil,
			ExpectedStatus: 400,
			LoginRes:       "",
			LoginErr:       nil,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrBadRequestBody.Error(),
			},
		},
		{
			Name: "body invalid",
			Body: map[string]interface{}{
				"email": "test@gmail.com",
			},
			ValidatorErr:   errors.New("invalid"),
			ExpectedStatus: 400,
			LoginRes:       "",
			LoginErr:       nil,
			ExpectedResult: map[string]interface{}{
				"message": "invalid",
			},
		},
		{
			Name: "user not found",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 400,
			LoginRes:       "",
			LoginErr:       utils.ErrNotFound,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrNotFound.Error(),
			},
		},
		{
			Name: "invalid password",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 400,
			LoginRes:       "",
			LoginErr:       utils.ErrInvalidPassword,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrInvalidPassword.Error(),
			},
		},
		{
			Name: "internal server error",
			Body: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			ValidatorErr:   nil,
			ExpectedStatus: 500,
			LoginRes:       "",
			LoginErr:       errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/login")

			//define mock
			mock1 := s.userServiceMock.On("Login").Return(v.LoginRes, v.LoginErr)
			mock2 := s.validaorMock.On("Validate").Return(v.ValidatorErr)

			err = s.userController.Login(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			mock1.Unset()
			mock2.Unset()
		})
	}
}

func (s *suiteUserController) TestGetUser() {
	// uuid statis in test
	varUUID := uuid.New()

	testCase := []struct {
		Name           string
		JwtRes         jwt.MapClaims
		ExpectedStatus int
		FindUserRes    *dto.UserResponse
		FindUserErr    error
		ExpectedResult map[string]interface{}
	}{
		{
			Name: "success",
			JwtRes: jwt.MapClaims{
				"user_id": varUUID.String(),
			},
			ExpectedStatus: 200,
			FindUserRes: &dto.UserResponse{
				ID: varUUID,
			},
			FindUserErr: nil,
			ExpectedResult: map[string]interface{}{
				"data": map[string]interface{}{
					"district_name": "",
					"email":         "",
					"id":            varUUID.String(),
					"name":          "",
					"phone":         "",
					"province_name": "",
					"regency_name":  "",
					"village_name":  "",
				},
				"message": "get user success",
			},
		},
		{
			Name: "internal server error",
			JwtRes: jwt.MapClaims{
				"user_id": varUUID.String(),
			},
			ExpectedStatus: 500,
			FindUserRes:    nil,
			FindUserErr:    errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/user")

			//define mock
			mock1 := s.userServiceMock.On("FindUser").Return(v.FindUserRes, v.FindUserErr)
			mock2 := s.JWTServiceMock.On("GetClaims").Return(v.JwtRes)

			err := s.userController.GetUser(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			mock1.Unset()
			mock2.Unset()
		})
	}
}

func (s *suiteUserController) TestGetAllUsers() {
	// uuid statis in test
	varUUID := uuid.New()

	testCase := []struct {
		Name            string
		JwtRes          jwt.MapClaims
		ExpectedStatus  int
		FindAllUsersRes dto.UsersResponse
		FindAllUsersErr error
		ExpectedResult  map[string]interface{}
	}{
		{
			Name: "success",
			JwtRes: jwt.MapClaims{
				"role_id": float64(3),
			},
			ExpectedStatus: 200,
			FindAllUsersRes: dto.UsersResponse{
				{
					ID: varUUID,
				},
			},
			FindAllUsersErr: nil,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"district_name": "",
						"email":         "",
						"id":            varUUID.String(),
						"name":          "",
						"phone":         "",
						"province_name": "",
						"regency_name":  "",
						"village_name":  "",
					},
				},
				"message": "success get users",
			},
		},
		{
			Name: "forbidden",
			JwtRes: jwt.MapClaims{
				"role_id": float64(1),
			},
			ExpectedStatus:  403,
			FindAllUsersRes: nil,
			FindAllUsersErr: nil,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrPermission.Error(),
			},
		},
		{
			Name: "internal server error",
			JwtRes: jwt.MapClaims{
				"role_id": float64(3),
			},
			ExpectedStatus:  500,
			FindAllUsersRes: nil,
			FindAllUsersErr: errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/user/all")

			//define mock
			mock1 := s.userServiceMock.On("FindAllUsers").Return(v.FindAllUsersRes, v.FindAllUsersErr)
			mock2 := s.JWTServiceMock.On("GetClaims").Return(v.JwtRes)

			err := s.userController.GetUsers(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			mock1.Unset()
			mock2.Unset()
		})
	}
}

func (s *suiteUserController) TestUpdateUser() {
	// uuid statis in test
	varUUID := uuid.New()

	testCase := []struct {
		Name           string
		Body           map[string]interface{}
		JwtRes         jwt.MapClaims
		ExpectedStatus int
		UpdateUserErr  error
		ExpectedResult map[string]interface{}
	}{
		{
			Name: "success",
			JwtRes: jwt.MapClaims{
				"user_id": varUUID.String(),
			},
			Body: map[string]interface{}{
				"name": "test",
			},
			ExpectedStatus: 200,
			UpdateUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"message": "success update user",
			},
		},
		{
			Name: "body type invalid",
			JwtRes: jwt.MapClaims{
				"user_id": varUUID.String(),
			},
			Body: map[string]interface{}{
				"name": 123,
			},
			UpdateUserErr:  nil,
			ExpectedStatus: 400,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrBadRequestBody.Error(),
			},
		},
		{
			Name: "error update",
			JwtRes: jwt.MapClaims{
				"user_id": varUUID.String(),
			},
			Body: map[string]interface{}{
				"name": "test",
			},
			ExpectedStatus: 400,
			UpdateUserErr:  errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			body, err := json.Marshal(v.Body)
			s.NoError(err)

			r := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/user/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(varUUID.String())

			// define mock
			mock1 := s.userServiceMock.On("UpdateUser").Return(v.UpdateUserErr)
			mock2 := s.JWTServiceMock.On("GetClaims").Return(v.JwtRes)

			err = s.userController.UpdateUser(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			// mock1.Unset()
			mock1.Unset()
			mock2.Unset()
		})
	}
}
func (s *suiteUserController) TestDeleteUser() {
	testCase := []struct {
		Name           string
		JwtRes         jwt.MapClaims
		ExpectedStatus int
		DeleteUserErr  error
		ExpectedResult map[string]interface{}
	}{
		{
			Name: "success",
			JwtRes: jwt.MapClaims{
				"role_id": float64(3),
			},
			ExpectedStatus: 200,
			DeleteUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"message": "success delete user",
			},
		},
		{
			Name: "forbidden",
			JwtRes: jwt.MapClaims{
				"role_id": float64(2),
			},
			ExpectedStatus: 403,
			DeleteUserErr:  nil,
			ExpectedResult: map[string]interface{}{
				"message": utils.ErrPermission.Error(),
			},
		},
		{
			Name: "delete error",
			JwtRes: jwt.MapClaims{
				"role_id": float64(3),
			},
			ExpectedStatus: 400,
			DeleteUserErr:  errors.New("err"),
			ExpectedResult: map[string]interface{}{
				"message": errors.New("err").Error(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodDelete, "/", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctx := s.echoNew.NewContext(r, w)
			ctx.SetPath("/user/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(uuid.New().String())

			// define mock
			mock1 := s.userServiceMock.On("DeleteUser").Return(v.DeleteUserErr)
			mock2 := s.JWTServiceMock.On("GetClaims").Return(v.JwtRes)

			err := s.userController.DeleteUser(ctx)
			s.NoError(err)

			controllerResult := map[string]interface{}{}
			err = json.NewDecoder(w.Result().Body).Decode(&controllerResult)
			s.NoError(err)

			s.Equal(v.ExpectedStatus, w.Result().StatusCode)
			s.Equal(v.ExpectedResult, controllerResult)

			// mock1.Unset()
			mock1.Unset()
			mock2.Unset()
		})
	}
}

func (s *suiteUserController) TestInitRoute() {
	group := s.echoNew.Group("/user")
	s.NotPanics(func() {
		s.userController.InitRoute(group, group)
	})
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(suiteUserController))
}
