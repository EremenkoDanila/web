package config

import (
	"os"
	"time"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"fmt"
	"strconv"
)


type Config struct {
	ServiceHost string
	ServicePort int
	MinioConfig MinioConfig
	JWT         JWTConfig

	Redis RedisConfig
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

type JWTConfig struct {
	Token         string        `mapstructure:"token"`          // секретный ключ
	ExpiresIn     time.Duration `mapstructure:"expires_in"`     // длительность жизни токена
	SigningMethod string        `mapstructure:"signing_method"` // "HS256"
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	// Устанавливаем значения по умолчанию для MinIO
	viper.SetDefault("minio.endpoint", "localhost:9000")
	viper.SetDefault("minio.access_key_id", "minioadmin")
	viper.SetDefault("minio.secret_access_key", "minioadmin")
	viper.SetDefault("minio.use_ssl", false)
	viper.SetDefault("minio.bucket_name", "software-images")

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}



	
	const (
		envRedisHost = "REDIS_HOST"
		envRedisPort = "REDIS_PORT"
		envRedisUser = "REDIS_USER"
		envRedisPass = "REDIS_PASSWORD"
	)
	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))
	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}
	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)




	log.Info("config parsed")

	return cfg, nil
}