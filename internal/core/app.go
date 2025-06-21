package core

import (
	"context"

	"github.com/hsingyingli/jarvis/pkg/config"
	"github.com/hsingyingli/jarvis/pkg/logger"
)

type App struct {
	config *config.Config
	logger logger.Logger
}

func NewApp(cfg *config.Config, log logger.Logger) (*App, error) {
	return &App{
		config: cfg,
		logger: log,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting Jarvis AI Assistant...")

	// TODO: Initialize services
	// TODO: Start HTTP server
	// TODO: Start WebSocket server

	<-ctx.Done()
	a.logger.Info("Application stopped")
	return nil
}