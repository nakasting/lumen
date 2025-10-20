package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Cannot read .env file %v", err)
	}

	cfg := &Config{
		Port: viper.GetString("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
