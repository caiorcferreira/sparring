package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
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

	err = e.Start(cfg.Port)
	return err
}
