package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func SetupTargets(cfg Config, e *echo.Echo) error {
	for _, target := range cfg.Targets {
		e.Add(
			target.Method,
			target.Path,
			makeTargetHandler(target),
			makeResponseTimeMiddleware(target),
		)
	}

	return nil
}

func makeResponseTimeMiddleware(target Target) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var resultErr error
			var wg sync.WaitGroup

			if target.ResponseTime > 0 {
				wg.Add(1)
				time.AfterFunc(target.ResponseTime, func() {
					resultErr = next(ctx)
					wg.Done()
				})
			}

			wg.Wait()

			return resultErr
		}
	}
}

func makeTargetHandler(target Target) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		body, err := mountTargetBody(target)
		if err != nil {
			ctx.Logger().Errorf("failed mounting target body for response: %s", err)
			return err
		}

		err = ctx.JSONBlob(target.StatusCode, body)
		if err != nil {
			ctx.Logger().Errorf("failed to write response back: %s", err)
			return err
		}

		return nil
	}
}

func mountTargetBody(target Target) ([]byte, error) {
	if target.Body.File != "" {
		return fetchFileContent(target.Body.File)
	}

	if target.Body.Value != "" {
		return []byte(target.Body.Value), nil
	}

	return nil, nil
}

func fetchFileContent(file string) ([]byte, error) {
	if !fileExists(file) {
		return nil, fmt.Errorf("body file %s not found", file)
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !info.IsDir()
}
