package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OllamaAPIURL string
	OllamaModel  string
	Port         string
	GinMode      string
	PostgresURL  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		OllamaAPIURL: os.Getenv("OLLAMA_API_URL"),
		OllamaModel:  getEnvOrDefault("OLLAMA_MODEL_NAME", "qwen2.5:3b"),
		Port:         getEnvOrDefault("PORT", "8080"),
		GinMode:      getEnvOrDefault("GIN_MODE", "release"),
		PostgresURL:  os.Getenv("POSTGRESQL_URL"),
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
