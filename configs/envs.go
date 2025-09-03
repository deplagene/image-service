package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	NsfwApiUrl    string
	MinioUser     string
	MinioPassword string
}

var Envs = mustLoad()

// todo: add logging
func mustLoad() *Config {
	const op = "configs.mustLoad"

	if err := godotenv.Load(); err != nil {
		return &Config{}
	}

	return &Config{
		NsfwApiUrl:    getEnv("NSFW_API_URL", "http://127.0.0.1:8000/v1/detect"),
		MinioUser:     getEnv("MINIO_ROOT_USER", "deplagene"),
		MinioPassword: getEnv("MINIO_ROOT_PASSWORD", "deplagene"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
