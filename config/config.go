package config

import (
	"github.com/spf13/viper"
)

const ServiceName = "cart"

type Config struct {
	Http *HttpConfig
	DB   *DBConfig
	Log  *LogConfig
}

func InitError(err error) {
	if err != nil {
		panic(err)
	}
}

func InitConfig() *Config {
	return &Config{
		Http: InitHttpConfig(),
		DB:   InitDBConfig(),
		Log:  InitLogConfig(),
	}
}

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("environment", "DEV")
	InitError(viper.BindEnv("environment", "ENVIRONMENT"))
}
