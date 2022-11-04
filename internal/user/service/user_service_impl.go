package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type userServiceImpl struct {
	repo     repository.UserRepository
	password utils.Password
}

// CreateAdmin implements UserService
func (u *userServiceImpl) CreateAdmin() error {
	var ctx context.Context
	users, err := u.repo.FindAllUsers(ctx)
	if err != nil {
		return err
	}
	if len(users) >= 1 {
		return nil
	}
	hashPassword, err := u.password.HashPassword(config.Cfg.DEFAULT_ADMIN_PASSWORD)
	if err != nil {
		return err
	}
	admin := model.User{
		ID:       uuid.New(),
		Name:     "admin",
		Email:    config.Cfg.DEFAULT_ADMIN_EMAIL,
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

func NewUserService(repository repository.UserRepository, password utils.Password) UserService {
	userService := &userServiceImpl{
		repo:     repository,
		password: password,
	}
	err := userService.CreateAdmin()
	if err != nil {
		panic(err)
	}

	return userService
}
