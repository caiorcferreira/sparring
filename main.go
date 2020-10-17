package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	if err := Start(); err != nil {
		log.Printf("[ERROR] server exited due to an error: %s", err)
		os.Exit(1)
	}
}

// Start runs an Sparring instance
func Start() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	e := echo.New()
	setLogLevel(cfg, e)

	e.GET("/health", func(ctx echo.Context) error {
		err := ctx.String(http.StatusOK, fmt.Sprintf("Healthy! :)\nConfig:\n%+v", cfg))
		if err != nil {
			log.Printf("[ERROR] failed to write response back: %s", err)
			return err
		}

		return nil
	})

	err = SetupTargets(cfg, e)
	if err != nil {
		return err
	}

	e.Use(middleware.Recover())
	if e.Logger.Level() <= log.INFO {
		e.Use(middleware.Logger())
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	go func() {
		if err := e.Start(addr); err != nil {
			e.Logger.Info("shutting down the server", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("graceful shutdown failed: %s", err)
	}

	return err
}

func setLogLevel(cfg Config, e *echo.Echo) {
	switch cfg.LogLevel {
	case "DEBUG", "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "INFO", "info":
		e.Logger.SetLevel(log.INFO)
	case "WARN", "warn":
		e.Logger.SetLevel(log.WARN)
	case "ERROR", "error":
		e.Logger.SetLevel(log.ERROR)
	}
}
