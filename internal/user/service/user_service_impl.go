package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type PasswordHashFunction interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type JWTService interface {
	GenerateToken(user *model.User) (string, error)
}
type userServiceImpl struct {
	repo       repository.UserRepository
	password   PasswordHashFunction
	jwtService JWTService
}

// Login implements UserService
func (u *userServiceImpl) Login(user dto.LoginRequest, ctx context.Context) (string, error) {
	userModel, err := u.repo.FindUserByEmail(user.Email, ctx)
	if err != nil {
		return "", err
	}
	if !u.password.CheckPasswordHash(user.Password, userModel.Password) {
		return "", utils.ErrInvalidPassword
	}
	token, err := u.jwtService.GenerateToken(userModel)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateAdmin implements UserService
func (u *userServiceImpl) CreateDefaultAdmin() error {
	var ctx context.Context
	users, err := u.repo.FindAllUsers(ctx)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return nil
	}
	hashPassword, err := u.password.HashPassword(constants.Default_password_admin)
	if err != nil {
		return err
	}
	admin := model.User{
		ID:       uuid.New(),
		Name:     "admin",
		Email:    constants.Default_email_admin,
		Password: hashPassword,
		RoleID:   3,
	}
	err = u.repo.CreateUser(&admin, ctx)
	return err
}

// CreateUser implements UserService
func (u *userServiceImpl) CreateUser(user dto.UserSignup, ctx context.Context) (uuid.UUID, error) {
	hashPassword, err := u.password.HashPassword(user.Password)
	if err != nil {
		return uuid.Nil, err
	}
	newID := uuid.New()
	userModel := user.ToModel()
	userModel.ID = newID
	userModel.Password = hashPassword
	userModel.RoleID = 1
	err = u.repo.CreateUser(userModel, ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return newID, nil
}

// FindAllUsers implements UserService
func (u *userServiceImpl) FindAllUsers(ctx context.Context) (dto.UsersResponse, error) {
	users, err := u.repo.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	var usersDto dto.UsersResponse
	usersDto.FromModel(users)
	return usersDto, nil
}

// FindUser implements UserService
func (u *userServiceImpl) FindUser(id string, ctx context.Context) (*dto.UserResponse, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	user, err := u.repo.FindUserByID(id, ctx)
	if err != nil {
		return nil, err
	}
	var userDto dto.UserResponse
	userDto.FromModel(user)
	return &userDto, nil
}

// UpdateUser implements UserService
func (u *userServiceImpl) UpdateUser(id string, user dto.UserUpdate, ctx context.Context) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrInvalidId
	}
	userModel := user.ToModel()
	userModel.ID = idUUID
	err = u.repo.UpdateUser(userModel, ctx)
	return err
}

// DeleteUser implements UserService
func (u *userServiceImpl) DeleteUser(id string, ctx context.Context) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrInvalidId
	}
	user := model.User{
		ID: idUUID,
	}
	err = u.repo.DeleteUser(&user, ctx)
	return err
}

func NewUserService(repository repository.UserRepository, password PasswordHashFunction, jwt JWTService) UserService {
	userService := &userServiceImpl{
		repo:       repository,
		password:   password,
		jwtService: jwt,
	}
	err := userService.CreateDefaultAdmin()
	if err != nil {
		panic(err)
	}

	return userService
}
