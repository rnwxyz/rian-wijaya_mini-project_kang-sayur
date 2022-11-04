package route

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pkgUserController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/controller"
	pkgUserRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	pkgUserService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Recover())
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}
	jwtService := utils.NewJWTService(config.Cfg.JWT_SECRET, 1*time.Hour)

	api := e.Group("/api")

	// version
	v1 := api.Group("/v1")
	auth := v1.Group("")
	auth.Use(middleware.JWT([]byte(config.Cfg.JWT_SECRET)))

	//init userController
	userRepository := pkgUserRepository.NewUserRepository(db)
	userService := pkgUserService.NewUserService(userRepository, utils.Password{}, jwtService)
	userController := pkgUserController.NewUserController(userService, jwtService)
	userController.InitRoute(v1, auth)
}
