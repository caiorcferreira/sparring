package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

func main() {
	if err := Start(); err != nil {
		log.Printf("[ERROR] server exited due to an error: %s", err)
		os.Exit(1)
	}
}

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
	if e.Logger.Level() == log.INFO {
		e.Use(middleware.Logger())
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	err = e.Start(addr)
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
