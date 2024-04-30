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

	// Gemini config
	GeminiApiKey   string `mapstructure:"GEMINI_API_KEY"`
	GeminiGenModel string `mapstructure:"GEMINI_GEN_MODEL"`
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
