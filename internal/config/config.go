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
	S3 struct {
		Bucket string `mapstructure:"bucket"`
		Region string `mapstructure:"region"`
	} `mapstructure:"s3"`
	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"int"`
	} `mapstructure:"redis"`
	Security struct {
		HMACEnabled   bool   `mapstructure:"hmac_enabled"`
		HMACSecretKey string `mapstructure:"hmac_secret_key"`
		RemoteFetch   struct {
			MaxDownloadSizeMB int `mapstructure:"max_download_size_mb"`
		} `mapstructure:"remote_fetch"`
	} `mapstructure:"security"`
}

func New() *Config {
	viper.SetDefault("http_server.port", 8080)
	viper.SetDefault("http_server.shutdown_timeout", "5s")
	viper.SetDefault("log.level", "info")

	viper.SetDefault("s3.bucket", "")
	viper.SetDefault("region", "eu-central-1")

	viper.SetDefault("redis.address", "localhost:6379")

	viper.SetDefault("security.hmac_secret_key", "")
	viper.SetDefault("security.hmac_enabled", true)

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

	if cfg.S3.Bucket == "" {
		log.Fatal("s3.bucket configuration is missing")
	}

	return &cfg
}
