package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("HTTP_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.HTTPPort != "9090" {
		t.Errorf("Expected HTTP_PORT to be 9090, got %s", cfg.HTTPPort)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LOG_LEVEL to be debug, got %s", cfg.LogLevel)
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.HTTPPort != "8080" {
		t.Errorf("Expected default HTTP_PORT to be 8080, got %s", cfg.HTTPPort)
	}

	if cfg.DBHost != "localhost" {
		t.Errorf("Expected default DB_HOST to be localhost, got %s", cfg.DBHost)
	}

	if cfg.DBPort != "5432" {
		t.Errorf("Expected default DB_PORT to be 5432, got %s", cfg.DBPort)
	}

	if cfg.RedisHost != "localhost" {
		t.Errorf("Expected default REDIS_HOST to be localhost, got %s", cfg.RedisHost)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default LOG_LEVEL to be info, got %s", cfg.LogLevel)
	}
}