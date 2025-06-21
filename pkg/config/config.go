package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server configuration
	HTTPPort string `json:"http_port"`
	HTTPHost string `json:"http_host"`

	// PostgreSQL configuration
	DatabaseURL string `json:"database_url"`
	DBHost      string `json:"db_host"`
	DBPort      string `json:"db_port"`
	DBName      string `json:"db_name"`
	DBUser      string `json:"db_user"`
	DBPassword  string `json:"db_password"`
	DBSSLMode   string `json:"db_ssl_mode"`

	// Redis configuration
	RedisURL      string `json:"redis_url"`
	RedisHost     string `json:"redis_host"`
	RedisPort     string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`

	// AI configuration
	OpenAIAPIKey string `json:"openai_api_key"`
	ModelName    string `json:"model_name"`

	// Logging configuration
	LogLevel string `json:"log_level"`

	// Security
	JWTSecret string `json:"jwt_secret"`
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		HTTPHost:      getEnv("HTTP_HOST", "localhost"),
		
		// PostgreSQL configuration
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "jarvis"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
		
		// Redis configuration
		RedisURL:      getEnv("REDIS_URL", ""),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		
		// AI configuration
		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
		ModelName:     getEnv("MODEL_NAME", "gpt-3.5-turbo"),
		
		// Other configuration
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}