package mock

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *model.User) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GetClaims(c *echo.Context) jwt.MapClaims {
	args := m.Called()

	return args.Get(0).(jwt.MapClaims)
}
