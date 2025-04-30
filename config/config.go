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
		Blockchain BlockchainConfig `yaml:"BlockchainConfig"`
		JWT        JWTConfig        `yaml:"JWT"`
		Kafka      KafkaConfig      `yaml:"Kafka"`
		Encryption EncryptionConfig `yaml:"Encryption"`
		Cors       CorsConfig       `yaml:"Cors"`
		GrpcServer GrpcServerConfig `yaml:"GrpcServer"`
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
	BlockchainConfig struct {
		GanacheURL             string `yaml:"GanacheURL"`
		VotechainAddress       string `yaml:"VotechainAddress" `
		VotechainBaseAddress   string `yaml:"VotechainBaseAddress"`
		KPUManagerAddress      string `yaml:"KPUManagerAddress"`
		VoterManagerAddress    string `yaml:"VoterManagerAddress"`
		ElectionManagerAddress string `yaml:"ElectionManagerAddress"`
	}
	JWTConfig struct {
		Secret string `yaml:"Secret" env:"JWT_SECRET"`
	}
	KafkaConfig struct {
		Producer KafkaProducerConfig `yaml:"Producer" env:"KAFKA_PRODUCER"`
		Topics   KafkaTopics         `yaml:"Topics" env:"KAFKA_TOPIC"`
	}
	KafkaProducerConfig struct {
		Brokers    []string `yaml:"Brokers" env:"KAFKA_BROKERS"`
		Idempotent bool     `yaml:"Idempotent" env:"KAFKA_IDEMPOTENT"`
		MaxAttempt int      `yaml:"MaxAttempt" env:"KAFKA_MAX_ATTEMPTS"`
	}

	KafkaTopics struct {
		MasterDataUser        KafkaTopicConfig `yaml:"MasterDataUser"`
		MasterDataVoter       KafkaTopicConfig `yaml:"MasterDataVoter"`
		MasterDataKPUProvinsi KafkaTopicConfig `yaml:"MasterDataKpuProvinsi"`
		MasterDataKPUKota     KafkaTopicConfig `yaml:"MasterDataKpuKota"`
	}
	GrpcServerConfig struct {
		Port uint `yaml:"Port"`
	}
	KafkaTopicConfig struct {
		Value        string `yaml:"Value" env:"KAFKA_TOPIC_VALUE"`
		ErrorHandler string `yaml:"ErrorHandler"`
		WithBackOff  bool   `yaml:"WithBackOff"`
	}
	EncryptionConfig struct {
		Key string `yaml:"Key" env:"ENCRYPTION_KEY"`
	}
	CorsConfig struct {
		AllowOrigins     string `yaml:"AllowOrigins"`
		AllowMethods     string `yaml:"AllowMethods"`
		AllowHeaders     string `yaml:"AllowHeaders"`
		AllowCredentials bool   `yaml:"AllowCredentials"`
		ExposeHeaders    string `yaml:"ExposeHeaders"`
		MaxAge           int    `yaml:"MaxAge"`
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
