package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	um "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	mm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/middleware/mock"
	pm "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/password/mock"
	"github.com/stretchr/testify/suite"
)

type suiteUserService struct {
	suite.Suite
	userRepositoryMock *um.UserRepositoryMock
	passwordMock       *pm.PasswordMock
	JWTServiceMock     *mm.MockJWTService
	userService        UserService
}

func newUserServiceMock(repository repository.UserRepository, password PasswordHashFunction, jwt JWTService) UserService {
	return &userServiceImpl{
		repo:       repository,
		password:   password,
		jwtService: jwt,
	}
}

func (s *suiteUserService) SetupTest() {
	s.userRepositoryMock = new(um.UserRepositoryMock)
	s.passwordMock = new(pm.PasswordMock)
	s.JWTServiceMock = new(mm.MockJWTService)
	s.userService = newUserServiceMock(s.userRepositoryMock, s.passwordMock, s.JWTServiceMock)
}

func (s *suiteUserService) TestCreateUser() {
	testCase := []struct {
		Name           string
		Body           dto.UserSignup
		ExpectedErr    error
		ExpectedResult uuid.UUID
		mockReturn     map[string]interface{}
	}{
		{
			Name: "success",
			Body: dto.UserSignup{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:    nil,
			ExpectedResult: uuid.New(),
			mockReturn: map[string]interface{}{
				"hashStr":    uuid.New().String(),
				"hashErr":    nil,
				"createUser": nil,
			},
		},
		{
			Name: "hash error",
			Body: dto.UserSignup{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:    errors.New("err"),
			ExpectedResult: uuid.Nil,
			mockReturn: map[string]interface{}{
				"hashStr":    uuid.New().String(),
				"hashErr":    errors.New("err"),
				"createUser": nil,
			},
		},
		{
			Name: "error create",
			Body: dto.UserSignup{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:    errors.New("err"),
			ExpectedResult: uuid.Nil,
			mockReturn: map[string]interface{}{
				"hashStr":    uuid.New().String(),
				"hashErr":    nil,
				"createUser": errors.New("err"),
			},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			hash := s.passwordMock.On("HashPassword").Return(v.mockReturn["hashStr"], v.mockReturn["hashErr"])
			create := s.userRepositoryMock.On("CreateUser").Return(v.mockReturn["createUser"])
			var ctx context.Context
			res, err := s.userService.CreateUser(v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(len(v.ExpectedResult.String()), len(res.String()))

			hash.Unset()
			create.Unset()
		})
	}
}

func (s *suiteUserService) TestLogin() {
	testCase := []struct {
		Name             string
		Body             dto.LoginRequest
		ExpectedErr      error
		ExpectedResult   string
		FindByEmailRes   *model.User
		FindByEmailErr   error
		CheckPassRes     bool
		GenerateTokenRes string
		GenerateTokenErr error
	}{
		{
			Name: "success",
			Body: dto.LoginRequest{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:      nil,
			ExpectedResult:   "test",
			FindByEmailRes:   &model.User{},
			FindByEmailErr:   nil,
			CheckPassRes:     true,
			GenerateTokenRes: "test",
			GenerateTokenErr: nil,
		},
		{
			Name: "error Find user",
			Body: dto.LoginRequest{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:      errors.New("err"),
			ExpectedResult:   "",
			FindByEmailRes:   &model.User{},
			FindByEmailErr:   errors.New("err"),
			CheckPassRes:     false,
			GenerateTokenRes: "",
			GenerateTokenErr: nil,
		},
		{
			Name: "Check password hash",
			Body: dto.LoginRequest{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:      customerrors.ErrInvalidPassword,
			ExpectedResult:   "",
			FindByEmailRes:   &model.User{},
			FindByEmailErr:   nil,
			CheckPassRes:     false,
			GenerateTokenRes: "",
			GenerateTokenErr: nil,
		},
		{
			Name: "error generate token",
			Body: dto.LoginRequest{
				Email:    "test@gmail.com",
				Password: "1234",
			},
			ExpectedErr:      errors.New("err"),
			ExpectedResult:   "",
			FindByEmailRes:   &model.User{},
			FindByEmailErr:   nil,
			CheckPassRes:     true,
			GenerateTokenRes: "",
			GenerateTokenErr: errors.New("err"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			fbe := s.userRepositoryMock.On("FindUserByEmail").Return(v.FindByEmailRes, v.FindByEmailErr)
			cpw := s.passwordMock.On("CheckPasswordHash").Return(v.CheckPassRes)
			gtk := s.JWTServiceMock.On("GenerateToken").Return(v.GenerateTokenRes, v.GenerateTokenErr)
			var ctx context.Context
			res, err := s.userService.Login(v.Body, ctx)
			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedResult, res)

			fbe.Unset()
			cpw.Unset()
			gtk.Unset()
		})
	}
}
func (s *suiteUserService) TestFindUser() {
	// uuid statis
	varUUID := uuid.New()

	testCase := []struct {
		Name            string
		Body            string
		ExpectedErr     error
		ExpectedResult  *dto.UserResponse
		FindUserByIdRes *model.User
		FindUserByIdErr error
	}{
		{
			Name:            "success",
			Body:            varUUID.String(),
			ExpectedErr:     nil,
			ExpectedResult:  &dto.UserResponse{ID: varUUID},
			FindUserByIdRes: &model.User{ID: varUUID},
			FindUserByIdErr: nil,
		},
		{
			Name:            "error id",
			Body:            "123",
			ExpectedErr:     customerrors.ErrInvalidId,
			ExpectedResult:  nil,
			FindUserByIdRes: nil,
			FindUserByIdErr: nil,
		},
		{
			Name:            "error user not found",
			Body:            uuid.New().String(),
			ExpectedErr:     errors.New("err"),
			ExpectedResult:  nil,
			FindUserByIdRes: nil,
			FindUserByIdErr: errors.New("err"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			fbi := s.userRepositoryMock.On("FindUserByID").Return(v.FindUserByIdRes, v.FindUserByIdErr)
			var ctx context.Context
			res, err := s.userService.FindUser(v.Body, ctx)
			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedResult, res)

			fbi.Unset()
		})
	}
}

func (s *suiteUserService) TestFindAllUsers() {
	testCase := []struct {
		Name            string
		ExpectedErr     error
		ExpectedResult  dto.UsersResponse
		FindAllUsersRes []model.User
		FindAllUsersErr error
	}{
		{
			Name:            "success",
			ExpectedErr:     nil,
			ExpectedResult:  dto.UsersResponse{{ID: uuid.New()}},
			FindAllUsersRes: []model.User{{ID: uuid.New()}},
			FindAllUsersErr: nil,
		},
		{
			Name:            "error when fild all users",
			ExpectedErr:     errors.New("err"),
			ExpectedResult:  nil,
			FindAllUsersRes: nil,
			FindAllUsersErr: errors.New("err"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			fau := s.userRepositoryMock.On("FindAllUsers").Return(v.FindAllUsersRes, v.FindAllUsersErr)
			var ctx context.Context
			res, err := s.userService.FindAllUsers(ctx)
			s.Equal(v.ExpectedErr, err)
			s.Equal(len(v.ExpectedResult), len(res))

			fau.Unset()
		})
	}
}

func (s *suiteUserService) TestUpdateUser() {
	testCase := []struct {
		Name          string
		Body          dto.UserUpdate
		Id            string
		ExpectedErr   error
		UpdateUserErr error
	}{
		{
			Name: "success",
			Body: dto.UserUpdate{
				Name: "test",
			},
			Id:            uuid.New().String(),
			ExpectedErr:   nil,
			UpdateUserErr: nil,
		},
		{
			Name: "invalid id",
			Body: dto.UserUpdate{
				Name: "test",
			},
			Id:            "123",
			ExpectedErr:   customerrors.ErrInvalidId,
			UpdateUserErr: nil,
		},
		{
			Name: "error update",
			Body: dto.UserUpdate{
				Name: "test",
			},
			Id:            uuid.New().String(),
			ExpectedErr:   errors.New("err"),
			UpdateUserErr: errors.New("err"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			upd := s.userRepositoryMock.On("UpdateUser").Return(v.UpdateUserErr)
			var ctx context.Context
			err := s.userService.UpdateUser(v.Id, v.Body, ctx)
			s.Equal(v.ExpectedErr, err)

			upd.Unset()
		})
	}
}
func (s *suiteUserService) TestDeleteUser() {
	testCase := []struct {
		Name          string
		Id            string
		ExpectedErr   error
		DeleteUserErr error
	}{
		{
			Name:          "success",
			Id:            uuid.New().String(),
			ExpectedErr:   nil,
			DeleteUserErr: nil,
		},
		{
			Name:          "invalid id",
			Id:            "123",
			ExpectedErr:   customerrors.ErrInvalidId,
			DeleteUserErr: nil,
		},
		{
			Name:          "error delete",
			Id:            uuid.New().String(),
			ExpectedErr:   errors.New("err"),
			DeleteUserErr: errors.New("err"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			upd := s.userRepositoryMock.On("DeleteUser").Return(v.DeleteUserErr)
			var ctx context.Context
			err := s.userService.DeleteUser(v.Id, ctx)
			s.Equal(v.ExpectedErr, err)

			upd.Unset()
		})
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(suiteUserService))
}
