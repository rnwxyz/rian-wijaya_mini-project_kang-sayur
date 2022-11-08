package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (b *UserServiceMock) CreateDefaultAdmin() error {
	args := b.Called()

	return args.Error(0)
}

func (b *UserServiceMock) CreateUser(user dto.UserSignup, ctx context.Context) (uuid.UUID, error) {
	args := b.Called()

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (b *UserServiceMock) UpdateUser(id string, user dto.UserUpdate, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}
func (b *UserServiceMock) FindAllUsers(ctx context.Context) (dto.UsersResponse, error) {
	args := b.Called()

	return args.Get(0).(dto.UsersResponse), args.Error(1)
}

func (b *UserServiceMock) Login(user dto.LoginRequest, ctx context.Context) (string, error) {
	args := b.Called()

	return args.Get(0).(string), args.Error(1)
}

func (b *UserServiceMock) FindUser(id string, ctx context.Context) (*dto.UserResponse, error) {
	args := b.Called()

	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (b *UserServiceMock) DeleteUser(id string, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}
