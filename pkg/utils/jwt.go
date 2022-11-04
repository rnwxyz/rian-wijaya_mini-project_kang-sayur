package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type JWTService struct {
	secretKey string
	exp       time.Duration
}

func NewJWTService(secretKey string, exp time.Duration) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		exp:       exp,
	}
}

func (j *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(j.exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) GetClaims(c *echo.Context) jwt.MapClaims {
	user := (*c).Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}
