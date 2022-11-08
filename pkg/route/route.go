package route

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pkgCheckpointController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/controller"
	pkgCheckpointRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/repository"
	pkgCheckpointService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/service"
	pkgItemController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/controller"
	pkgItemRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	pkgItemService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/service"
	pkgOrderController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/controller"
	pkgOrderRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	pkgOrderService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/service"
	pkgRegionController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/controller"
	pkgRegionRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/repository"
	pkgRegionService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/service"
	pkgTransactionController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/controller"
	pkgTransactionRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/repository"
	pkgTransactionService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/transaction/service"
	pkgUserController "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/controller"
	pkgUserRepository "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	pkgUserService "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/service"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	importcsv "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/import_csv"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/payment"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

func InitGlobalRoute(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Recover())
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}
	jwtService := utils.NewJWTService(config.Cfg.JWT_SECRET, constants.ExpToken)

	api := e.Group("/api")

	// version
	v1 := api.Group("/v1")
	auth := v1.Group("")
	auth.Use(middleware.JWT([]byte(config.Cfg.JWT_SECRET)))

	//init user controller
	userRepository := pkgUserRepository.NewUserRepository(db)
	userService := pkgUserService.NewUserService(userRepository, utils.Password{}, jwtService)
	userController := pkgUserController.NewUserController(userService, jwtService)
	userController.InitRoute(v1, auth)

	//init region controller
	importCsvService := importcsv.NewImportCsv()
	regionRepository := pkgRegionRepository.NewRegionRepository(db)
	regionService := pkgRegionService.NewRegionService(regionRepository, importCsvService)
	regionController := pkgRegionController.NewRegionController(regionService)
	regionController.InitRoute(auth)

	// init checkpoint controller
	checkpointRepository := pkgCheckpointRepository.NewCheckpointRepository(db)
	checkpointService := pkgCheckpointService.NewCheckpointService(checkpointRepository)
	checkpointController := pkgCheckpointController.NewCheckpointController(checkpointService, jwtService)
	checkpointController.InitRoute(auth)

	// init item controller
	itemRepository := pkgItemRepository.NewItemRepository(db)
	itemService := pkgItemService.NewItemService(itemRepository)
	itemController := pkgItemController.NewItemController(itemService, jwtService)
	itemController.InitRoute(auth)

	// init order controller
	orderRepository := pkgOrderRepository.NewOrderRepository(db)
	orderService := pkgOrderService.NewOrderService(orderRepository, itemRepository, payment.Midtrans{}, userRepository)
	orderController := pkgOrderController.NewOrderController(orderService, jwtService)
	orderController.InitRoute(auth)

	// init transaction controller
	transactionRepository := pkgTransactionRepository.NewTransactionRepository(db)
	transactionService := pkgTransactionService.NewTransactionService(transactionRepository, orderRepository)
	transactionController := pkgTransactionController.NewTransactionController(transactionService, jwtService)
	transactionController.InitRoute(v1, auth)
}
