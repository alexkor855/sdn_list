package config

import "github.com/spf13/viper"

type HttpConfig struct {
	Host string
	Port string
}

func InitHttpConfig() *HttpConfig {
	httpConfig := HttpConfig{
		Host: viper.GetString("http.host"),
		Port: viper.GetString("http.port"),
	}
	return &httpConfig
}

func init() {
	viper.SetDefault("http.host", ServiceName)
	InitError(viper.BindEnv("http.host", "HTTP_HOST"))
	InitError(viper.BindEnv("http.port", "HTTP_PORT"))
}
