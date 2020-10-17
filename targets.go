package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SetupTargets(cfg Config, e *echo.Echo) error {
	for _, target := range cfg.Targets {
		e.Add(target.Method, target.Path, func(ctx echo.Context) error {
			err := ctx.JSONBlob(http.StatusOK, []byte(target.Body))
			if err != nil {
				e.Logger.Errorf("failed to write response back: %s", err)
				return err
			}

			return nil
		})
	}

	return nil
}
