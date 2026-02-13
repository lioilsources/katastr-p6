package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	RedisURL    string
	CUZKAPIKey  string
	CUZKBaseURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		CUZKAPIKey:  getEnv("CUZK_API_KEY", ""),
		CUZKBaseURL: getEnv("CUZK_BASE_URL", "https://api-kn.cuzk.gov.cz/api/v1"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
