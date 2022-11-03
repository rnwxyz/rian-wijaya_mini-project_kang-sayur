package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
)

type UserService interface {
	CreateUser(user dto.UserRequest, ctx context.Context) (uuid.UUID, error)
	UpdateUser(id string, user dto.UserRequest, ctx context.Context) error
	FindAllUsers(ctx context.Context) (dto.UsersResponse, error)
}
