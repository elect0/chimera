package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpSever struct {
		Port            int           `mapstructure:"port"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	} `mapstructure:"http_server"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}

func New() *Config {
	viper.SetDefault("http_server.port", 8080)
	viper.SetDefault("http_server.shutdown_timeout", "5s")
	viper.SetDefault("log.level", "info")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found; using defaults")
		} else {
			log.Fatalf("fatal error config file: %s", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &cfg
}
