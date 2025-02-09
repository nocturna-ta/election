package config

import (
	"github.com/nocturna-ta/golib/config"
	"github.com/nocturna-ta/golib/log"
	"time"
)

type (
	MainConfig struct {
		Server     ServerConfig     `yaml:"Server"`
		API        APIConfig        `yaml:"API"`
		Database   DBConfig         `yaml:"Database"`
		Encryption EncryptionConfig `yaml:"Encryption"`
	}
	ServerConfig struct {
		Port         uint          `yaml:"Port" env:"SERVER_PORT"`
		WriteTimeout time.Duration `yaml:"WriteTimeout" env:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `yaml:"ReadTimeout" env:"SERVER_READ_TIMEOUT"`
	}
	APIConfig struct {
		BasePath      string        `yaml:"BasePath" env:"API_BASE_PATH"`
		APITimeout    time.Duration `yaml:"APITimeout" env:"API_TIMEOUT"`
		EnableSwagger bool          `yaml:"EnableSwagger" env:"ENABLE_SWAGGER" default:"false"`
	}
	DBConfig struct {
		SlaveDSN        string `yaml:"SlaveDSN" env:"DB_SLAVE_DSN"`
		MasterDSN       string `yaml:"MasterDSN" env:"DB_MASTER_DSN"`
		RetryInterval   int    `yaml:"RetryInterval" env:"DB_RETRY_INTERVAL"`
		MaxIdleConn     int    `yaml:"MaxIdleConn" env:"DB_MAX_IDLE_CONN"`
		MaxConn         int    `yaml:"MaxConn" env:"DB_MAX_CONN"`
		ConnMaxLifetime string `yaml:"ConnMaxLifetime" env:"DB_CONN_MAX_LIFETIME"`
	}
	EncryptionConfig struct {
		Key string `yaml:"Key" env:"ENCRYPTION_KEY"`
	}
)

func ReadConfig(cfg any, configLocation string) {
	if configLocation == "" {
		configLocation = "file://config/files/config.yaml"
	}

	if err := config.ReadConfig(cfg, configLocation, true); err != nil {
		log.WithFields(log.Fields{
			"error":           err,
			"config-location": configLocation,
		}).Fatal("Failed to read config")
	}
}
