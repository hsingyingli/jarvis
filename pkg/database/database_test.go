package database

import (
	"testing"
	"github.com/hsingyingli/jarvis/pkg/config"
)

func TestNew(t *testing.T) {
	// Skip this test if no PostgreSQL/Redis available
	t.Skip("Skipping database test - requires PostgreSQL and Redis running")
	
	cfg := &config.Config{
		DatabaseURL: "postgresql://postgres:password@localhost:5432/jarvis_test?sslmode=disable",
		RedisURL:    "redis://localhost:6379/1",
	}

	db, err := New(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.Close()

	if db.Postgres == nil {
		t.Error("Expected PostgreSQL connection, got nil")
	}

	if db.Redis == nil {
		t.Error("Expected Redis connection, got nil")
	}
}

func TestNewWithIndividualConfig(t *testing.T) {
	// Skip this test if no PostgreSQL/Redis available
	t.Skip("Skipping database test - requires PostgreSQL and Redis running")

	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432", 
		DBName:     "jarvis_test",
		DBUser:     "postgres",
		DBPassword: "password",
		DBSSLMode:  "disable",
		RedisHost:  "localhost",
		RedisPort:  "6379",
		RedisDB:    1,
	}

	db, err := New(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.Close()

	if db.Postgres == nil {
		t.Error("Expected PostgreSQL connection, got nil")
	}

	if db.Redis == nil {
		t.Error("Expected Redis connection, got nil")
	}
}

func TestInvalidConfig(t *testing.T) {
	cfg := &config.Config{
		DatabaseURL: "invalid://connection/string",
		RedisURL:    "invalid://redis/url",
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error for invalid configuration")
	}
}