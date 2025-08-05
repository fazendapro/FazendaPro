package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBHost    string
	DBPort    string
	User      string
	Password  string
	Name      string
	CORS      CORSConfig
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("WARNING: Não foi possível carregar .env: %v\n", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("WARNING: Não foi possível carregar %s: %v\n", envFile, err)
	}

	fmt.Printf("DEBUG: ENV=%s, DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s\n",
		os.Getenv("ENV"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"))

	return &Config{
		Port:      getEnvWithDefault("PORT", "8080"),
		JWTSecret: getEnvWithDefault("JWT_SECRET", "dev-secret-key"),
		DBHost:    getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:    getEnvWithDefault("DB_PORT", "5432"),
		User:      getEnvWithDefault("DB_USER", "fazendapro_user"),
		Password:  getEnvWithDefault("DB_PASSWORD", "fazendapro_password"),
		Name:      getEnvWithDefault("DB_NAME", "fazendapro"),
		CORS:      loadCORSConfig(),
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadCORSConfig() CORSConfig {
	corsConfig := CORSConfig{
		AllowedOrigins:   splitEnvVar(getEnvWithDefault("CORS_ALLOWED_ORIGINS", "*")),
		AllowedMethods:   splitEnvVar(getEnvWithDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")),
		AllowedHeaders:   splitEnvVar(getEnvWithDefault("CORS_ALLOWED_HEADERS", "Content-Type,Authorization")),
		ExposedHeaders:   splitEnvVar(getEnvWithDefault("CORS_EXPOSED_HEADERS", "")),
		AllowCredentials: getEnvWithDefault("CORS_ALLOW_CREDENTIALS", "true") == "true",
		MaxAge:           parseInt(getEnvWithDefault("CORS_MAX_AGE", "86400")),
	}

	return corsConfig
}

func splitEnvVar(value string) []string {
	if value == "" {
		return []string{}
	}

	var result []string
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseInt(value string) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return 0
}
