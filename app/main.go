package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/database"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/route"
)

func main() {
	config.InitConfig()

	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	if err := database.MigrateDB(db); err != nil {
		panic(err)
	}
	e := echo.New()
	route.InitGlobalRoute(e, db)
	e.Logger.Fatal(e.Start(":" + config.Cfg.API_PORT))
}
