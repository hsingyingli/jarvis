package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hsingyingli/jarvis/internal/core"
	"github.com/hsingyingli/jarvis/pkg/config"
	"github.com/hsingyingli/jarvis/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize core application
	app, err := core.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize application:", err)
	}

	// Start application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		cancel()
	}()

	if err := app.Run(ctx); err != nil {
		logger.Fatal("Application error:", err)
	}
}