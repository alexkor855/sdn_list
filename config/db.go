package config

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	ConnString string
}

func InitDBConfig() *DBConfig {	
	dbConfig := DBConfig{
		ConnString: viper.GetString("db.conn_string"),
	}
	return &dbConfig
}

func init() {
	InitError(viper.BindEnv("db.conn_string", "DATABASE_URL"))
}
