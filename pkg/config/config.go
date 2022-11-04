package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	API_PORT               string
	DB_ADDRESS             string
	DB_USERNAME            string
	DB_PASSWORD            string
	DB_NAME                string
	DEFAULT_ADMIN_EMAIL    string
	DEFAULT_ADMIN_PASSWORD string
}

var Cfg *Config

func InitConfig() {
	cfg := &Config{}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	viper.Unmarshal(cfg)

	Cfg = cfg
}
