package controller

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pkgUserController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/controller"
	pkgUserRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	pkgUserService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

func InitController(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Recover())
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}

	api := e.Group("/api")

	// version 1
	v1 := api.Group("/v1")

	//init userController
	userRepository := pkgUserRepository.NewUserRepository(db)
	userService := pkgUserService.NewUserService(userRepository, utils.Password{})
	userController := pkgUserController.NewUserController(userService)
	userController.InitRoute(v1)
}
