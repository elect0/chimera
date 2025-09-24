package config

import "time"

type Config struct {
	HttpSever struct {
		Port            int           `mapstructure:"port"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	}
	Log struct {
		Level string `mapstructure:"level"`
	}
}

func New() *Config {
}
