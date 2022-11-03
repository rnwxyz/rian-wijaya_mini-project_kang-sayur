package main

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/database"
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
}
