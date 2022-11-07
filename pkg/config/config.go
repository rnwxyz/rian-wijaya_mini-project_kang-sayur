package config

import "os"

type Config struct {
	API_PORT               string
	DB_ADDRESS             string
	DB_USERNAME            string
	DB_PASSWORD            string
	DB_NAME                string
	DEFAULT_ADMIN_EMAIL    string
	DEFAULT_ADMIN_PASSWORD string
	JWT_SECRET             string
	TIME_LOCATION          string
	ORDER_SECRET           string
	MIDTRANS_SERVER_KEY    string
	DNS                    string
}

var Cfg *Config

func InitConfig() {
	cfg := &Config{}

	cfg.DB_ADDRESS = os.Getenv("DB_ADDRESS")
	cfg.API_PORT = os.Getenv("API")
	cfg.DB_USERNAME = os.Getenv("DB_USERNAME")
	cfg.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	cfg.DB_NAME = os.Getenv("DB_NAME")
	cfg.DEFAULT_ADMIN_EMAIL = os.Getenv("DEFAULT_ADMIN_EMAIL")
	cfg.DEFAULT_ADMIN_PASSWORD = os.Getenv("DEFAULT_ADMIN_PASSWORD")
	cfg.JWT_SECRET = os.Getenv("JWT_SECRET")
	cfg.TIME_LOCATION = os.Getenv("TIME_LOCATION")
	cfg.ORDER_SECRET = os.Getenv("ORDER_SECRET")
	cfg.MIDTRANS_SERVER_KEY = os.Getenv("MIDTRANS_SERVER_KEY")
	cfg.DNS = os.Getenv("DNS")

	Cfg = cfg
}
