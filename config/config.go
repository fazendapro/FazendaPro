package config

import (
	"fmt"
	"net/url"
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

	var dbHost, dbPort, dbUser, dbPassword, dbName string

	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		fmt.Printf("DEBUG: Usando DB_URL para configuração do banco\n")
		var err error
		dbHost, dbPort, dbUser, dbPassword, dbName, err = parseDBURL(dbURL)
		if err != nil {
			return nil, fmt.Errorf("erro ao processar DB_URL: %v", err)
		}
	} else {
		fmt.Printf("DEBUG: Usando configuração manual do banco\n")
		dbHost = getEnvWithDefault("DB_HOST", "localhost")
		dbPort = getEnvWithDefault("DB_PORT", "5432")
		dbUser = getEnvWithDefault("DB_USER", "fazendapro_user")
		dbPassword = getEnvWithDefault("DB_PASSWORD", "fazendapro_password")
		dbName = getEnvWithDefault("DB_NAME", "fazendapro")
	}

	fmt.Printf("DEBUG: ENV=%s, DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s\n",
		os.Getenv("ENV"),
		dbHost,
		dbPort,
		dbUser,
		dbName)

	return &Config{
		Port:      getEnvWithDefault("PORT", "8080"),
		JWTSecret: getEnvWithDefault("JWT_SECRET", "dev-secret-key"),
		DBHost:    dbHost,
		DBPort:    dbPort,
		User:      dbUser,
		Password:  dbPassword,
		Name:      dbName,
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

func parseDBURL(dbURL string) (host, port, user, password, dbName string, err error) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("erro ao fazer parse da DB_URL: %v", err)
	}

	host = parsedURL.Hostname()
	port = parsedURL.Port()
	if port == "" {
		port = "5432"
	}

	if parsedURL.User != nil {
		user = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	dbName = strings.TrimPrefix(parsedURL.Path, "/")

	return host, port, user, password, dbName, nil
}
