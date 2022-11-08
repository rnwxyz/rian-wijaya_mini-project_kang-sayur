package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (b *UserRepositoryMock) CreateUser(user *model.User, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}

func (b *UserRepositoryMock) FindAllUsers(ctx context.Context) ([]model.User, error) {
	args := b.Called()

	return args.Get(0).([]model.User), args.Error(1)
}

func (b *UserRepositoryMock) FindUserByEmail(email string, ctx context.Context) (*model.User, error) {
	args := b.Called()

	return args.Get(0).(*model.User), args.Error(1)
}

func (b *UserRepositoryMock) FindUserByID(id string, ctx context.Context) (*model.User, error) {
	args := b.Called()

	return args.Get(0).(*model.User), args.Error(1)
}

func (b *UserRepositoryMock) DeleteUser(user *model.User, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}

func (b *UserRepositoryMock) UpdateUser(user *model.User, ctx context.Context) error {
	args := b.Called()

	return args.Error(0)
}

func (b *UserRepositoryMock) InitRole() error {
	args := b.Called()

	return args.Error(0)
}
