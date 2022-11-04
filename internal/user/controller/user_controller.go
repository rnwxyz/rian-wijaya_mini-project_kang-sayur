package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *userController {
	return &userController{
		service: service,
	}
}

func (u *userController) InitRoute(e *echo.Group) {
	e.POST("/signup", u.SignUp)
	e.GET("/user", u.GetUsers)
	e.PUT("/user/:id", u.UpdateUser)
	e.DELETE("/user/:id", u.DeleteUser)
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

func (u *userController) GetUsers(c echo.Context) error {
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
	var user dto.UserUpdate
	id := c.Param("id")
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": utils.ErrBadRequestBody.Error()})
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	err := u.service.UpdateUser(id, user, c.Request().Context())
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
