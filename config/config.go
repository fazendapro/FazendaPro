package config

import (
	"fmt"
	"os"

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
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("WARNING: Não foi possível carregar .env: %v\n", err)
	}

	// Determinar ambiente
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Carregar arquivo de ambiente específico se existir
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("WARNING: Não foi possível carregar %s: %v\n", envFile, err)
	}

	// Debug: imprimir configurações carregadas
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
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
