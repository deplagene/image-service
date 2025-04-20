package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	S3Url              string
	S3User             string
	S3Password         string
	BrokerUser         string
	BrokerPassword     string
	BrokerHost         string
	DatabaseConnString string
	ApiHost            string
	ApiPort            string
	Env                string
}

var Envs = initConfig()

func initConfig() *Config {
	if err := godotenv.Load(); err != nil {
		return &Config{}
	}

	return &Config{
		S3Url:              getEnv("MINIO_ROOT_URL", "localhost:9000"),
		S3User:             getEnv("MINIO_ROOT_USER", "minio"),
		S3Password:         getEnv("MINIO_ROOT_PASSWORD", "minio"),
		BrokerUser:         getEnv("RABBITMQ_USERNAME", "guest"),
		BrokerPassword:     getEnv("RABBITMQ_PASSWORD", "guest"),
		BrokerHost:         getEnv("RABBITMQ_HOST", "localhost:5672"),
		DatabaseConnString: getEnv("DB_CONNECTION_STRING", "postgresql://postgres:postgres@localhost:5432/images_service"),
		ApiHost:            getEnv("API_HOST", "localhost"),
		ApiPort:            getEnv("API_PORT", "8080"),
		Env:                getEnv("ENVIRONMENT", "default"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if val, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
