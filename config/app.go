package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// App configuration
	LogLevel string `mapstructure:"LOG_LEVEL"`
	Env      string `mapstructure:"ENV"`
	Port     string `mapstructure:"PORT"`

	// Client config
	ClientUrl string `mapstructure:"CLIENT_URL"`

	// DB config
	DbHost    string `mapstructure:"DB_HOST"`
	DbPort    string `mapstructure:"DB_PORT"`
	DbName    string `mapstructure:"DB_NAME"`
	DbUser    string `mapstructure:"DB_USER"`
	DbPwd     string `mapstructure:"DB_PWD"`
	DbSslmode string `mapstructure:"DB_SSLMODE"`
}

func LoadConfig() (config Config) {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to load config: ", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("failed to load config: ", err)
	}

	configureLogger(config)

	return config
}
