package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func SetupTargets(cfg Config, e *echo.Echo) error {
	for _, target := range cfg.Targets {
		e.Add(target.Method, target.Path, makeTargetHandler(target))
	}

	return nil
}

func makeTargetHandler(target Target) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if target.ResponseTime != 0 {
			time.Sleep(target.ResponseTime)
		}

		err := ctx.JSONBlob(http.StatusOK, []byte(target.Body))
		if err != nil {
			ctx.Logger().Errorf("failed to write response back: %s", err)
			return err
		}

		return nil
	}
}
