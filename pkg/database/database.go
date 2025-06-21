package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/hsingyingli/jarvis/pkg/config"
)

type DB struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

func New(cfg *config.Config) (*DB, error) {
	// Connect to PostgreSQL
	var pgConnString string
	if cfg.DatabaseURL != "" {
		pgConnString = cfg.DatabaseURL
	} else {
		pgConnString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	}

	pgPool, err := pgxpool.New(context.Background(), pgConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test PostgreSQL connection
	if err := pgPool.Ping(context.Background()); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Connect to Redis
	var redisOpt *redis.Options
	if cfg.RedisURL != "" {
		redisOpt, err = redis.ParseURL(cfg.RedisURL)
		if err != nil {
			pgPool.Close()
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
	} else {
		redisOpt = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}
	}

	redisClient := redis.NewClient(redisOpt)

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &DB{
		Postgres: pgPool,
		Redis:    redisClient,
	}, nil
}

func (db *DB) Close() error {
	var pgErr, redisErr error
	
	if db.Postgres != nil {
		db.Postgres.Close()
	}
	
	if db.Redis != nil {
		redisErr = db.Redis.Close()
	}
	
	if pgErr != nil {
		return fmt.Errorf("PostgreSQL close error: %w", pgErr)
	}
	if redisErr != nil {
		return fmt.Errorf("Redis close error: %w", redisErr)
	}
	
	return nil
}