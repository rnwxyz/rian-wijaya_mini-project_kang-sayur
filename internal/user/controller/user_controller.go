package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type (
	JWTService interface {
		GetClaims(c *echo.Context) jwt.MapClaims
	}
	userController struct {
		service    service.UserService
		jwtService JWTService
	}
)

func NewUserController(service service.UserService, jwt JWTService) *userController {
	return &userController{
		service:    service,
		jwtService: jwt,
	}
}

func (u *userController) InitRoute(api *echo.Group, auth *echo.Group) {
	api.POST("/signup", u.SignUp)
	api.POST("/login", u.Login)
	auth.GET("/user", u.GetUsers)
	auth.GET("/user/:id", u.GetUser)
	auth.PUT("/user/:id", u.UpdateUser)
	auth.DELETE("/user/:id", u.DeleteUser)
}

func (u *userController) SignUp(c echo.Context) error {
	var user dto.UserSignup
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	id, err := u.service.CreateUser(user, c.Request().Context())
	if err != nil {
		if err == utils.ErrEmailAlredyExist {
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
			"message": utils.ErrBadRequestBody.Error(),
		})
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	token, err := u.service.Login(user, c.Request().Context())
	if err != nil {
		if err == utils.ErrNotFound || err == utils.ErrInvalidPassword {
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
	paramId := c.Param("id")
	if userId != paramId {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	user, err := u.service.FindUser(paramId, c.Request().Context())
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
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
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
	paramId := c.Param("id")
	if userId != paramId {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
		})
	}
	var user dto.UserUpdate
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	err := u.service.UpdateUser(paramId, user, c.Request().Context())
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
	if role < 3 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": utils.ErrPermission.Error(),
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
