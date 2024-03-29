package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}
type checkpointController struct {
	service    service.CheckpointService
	jwtService JWTService
}

func NewCheckpointController(service service.CheckpointService, jwt JWTService) *checkpointController {
	return &checkpointController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *checkpointController) InitRoute(auth *echo.Group) {
	checkpoints := auth.Group("/checkpoints")
	checkpoints.POST("", u.CreateCheckpoint)
	checkpoints.GET("", u.GetCheckpoints)
	checkpoints.GET("/profile", u.GetCheckpointByUser)
}

func (u *checkpointController) CreateCheckpoint(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role != constants.Role_admin {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": customerrors.ErrPermission.Error(),
		})
	}
	var checkpointBody dto.CheckpointRequest
	if err := c.Bind(&checkpointBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(checkpointBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error()})
	}
	id, err := u.service.CreateCheckpoint(checkpointBody, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrBadRequestBody {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "new checkpoint success created",
		"id":      id,
	})
}

func (u *checkpointController) GetCheckpoints(c echo.Context) error {
	checkpoints, err := u.service.FindCheckpoints(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get checkpoint success",
		"data":    checkpoints,
	})
}

func (u *checkpointController) GetCheckpointByUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)
	checkpoints, err := u.service.FindCheckpointsByUser(userId, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrCheckpointNotCovered {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get checkpoint success",
		"data":    checkpoints,
	})
}
