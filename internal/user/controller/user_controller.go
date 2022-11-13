package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type JWTService interface {
	GetClaims(c *echo.Context) jwt.MapClaims
}

type userController struct {
	service    service.UserService
	jwtService JWTService
}

func NewUserController(service service.UserService, jwt JWTService) *userController {
	return &userController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *userController) InitRoute(api *echo.Group, auth *echo.Group) {
	api.POST("/signup", u.SignUp)
	api.POST("/login", u.Login)

	users := auth.Group("/users")
	users.GET("", u.GetUsers)
	users.PUT("", u.UpdateUser)
	users.DELETE("/:id", u.DeleteUser)
	users.GET("/profile", u.GetUser)
}

func (u *userController) SignUp(c echo.Context) error {
	var user dto.UserSignup
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	id, err := u.service.CreateUser(user, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrEmailAlredyExist {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "sign up success",
		"id":      id,
	})
}

func (u *userController) Login(c echo.Context) error {
	var user dto.LoginRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrBadRequestBody.Error(),
		})
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	token, err := u.service.Login(user, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrNotFound || err == customerrors.ErrInvalidPassword {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "login success",
		"token":   token,
	})
}

func (u *userController) GetUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)
	user, err := u.service.FindUser(userId, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get user success",
		"data":    user,
	})
}

func (u *userController) GetUsers(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role != constants.Role_admin {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": customerrors.ErrPermission.Error(),
		})
	}
	users, err := u.service.FindAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get users",
		"data":    users,
	})
}

func (u *userController) UpdateUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userId := claims["user_id"].(string)
	var user dto.UserUpdate
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": customerrors.ErrBadRequestBody.Error()})
	}
	err := u.service.UpdateUser(userId, user, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update user",
	})
}

func (u *userController) DeleteUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role_id"].(float64)
	if role != constants.Role_admin {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": customerrors.ErrPermission.Error(),
		})
	}
	id := c.Param("id")
	err := u.service.DeleteUser(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success delete user",
	})
}
