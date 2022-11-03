package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type UserRepository interface {
	CreateUser(user *model.User, ctx context.Context) error
	UpdateUser(user *model.User, ctx context.Context) error
	DeleteUser(user *model.User, ctx context.Context) error
	FindAllUsers(ctx context.Context) ([]model.User, error)
	FindSingleUser(user *model.User, ctx context.Context) error
	FindMultipleUsers(user *model.User, ctx context.Context) ([]model.User, error)
}
