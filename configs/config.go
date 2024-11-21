package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Mode          string `envconfig:"MODE" default:"debug"`
	DatabaseTrx   DatabaseTrx
	DatabaseParam DatabaseParam
	DatabaseMerchant DatabaseMerchant
	ServerPort    int `envconfig:"SERVER_PORT" default:"88"`
	TimeoutTrx	  int `envconfig:"TIMEOUT_TRX" default:"30"` 
}

type DatabaseTrx struct {
	Host     string `envconfig:"DATABASE_TRX_HOST" required:"true"`
	Port     int    `envconfig:"DATABASE_TRX_PORT" required:"true"`
	User     string `envconfig:"DATABASE_TRX_USER" required:"true"`
	Password string `envconfig:"DATABASE_TRX_PASSWORD" required:"true"`
	Name     string `envconfig:"DATABASE_TRX_NAME" required:"true"`
}

type DatabaseParam struct {
	Host     string `envconfig:"DATABASE_PARAM_HOST" required:"true"`
	Port     int    `envconfig:"DATABASE_PARAM_PORT" required:"true"`
	User     string `envconfig:"DATABASE_PARAM_USER" required:"true"`
	Password string `envconfig:"DATABASE_PARAM_PASSWORD" required:"true"`
	Name     string `envconfig:"DATABASE_PARAM_NAME" required:"true"`
}

type DatabaseMerchant struct {
	Host     string `envconfig:"DATABASE_MERCHANT_HOST"`
	Port     int    `envconfig:"DATABASE_MERCHANT_PORT"`
	User     string `envconfig:"DATABASE_MERCHANT_USER"`
	Password string `envconfig:"DATABASE_MERCHANT_PASSWORD"`
	Name     string `envconfig:"DATABASE_MERCHANT_NAME"`
}

func NewParsedConfig() (Config, error) {
	_ = godotenv.Load(".env")
	cnf := Config{}
	err := envconfig.Process("", &cnf)
	return cnf, err
}