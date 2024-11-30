package app

import (
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/logger"
)

// App struct of the application.
type App struct {
	cfg *config.Config
}

// New creates new application.
func New() (*App, error) {
	app := &App{}

	// Configuration initialization
	err := app.initConfig()
	if err != nil {
		return nil, err
	}

	// Logger initialization
	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	return app, nil
}

// initConfig initializes application configuration.
func (a *App) initConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	a.cfg = cfg

	return nil
}

// initLogger initializes logger.
func (a *App) initLogger() error {
	err := logger.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// Run runs the application.
func (a *App) Run() error {
	return nil
}
