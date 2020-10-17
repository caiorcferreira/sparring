package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var cache sync.Map

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
			ctx.Logger().Errorf("failed to mount target response body: %s", err)
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
	cacheKey := fmt.Sprintf("%s %s", target.Method, target.Path)
	value, found := cache.Load(cacheKey)
	if found {
		b, ok := value.([]byte)
		if ok {
			return b, nil
		}
	}

	var body []byte
	if target.Body.File != "" {
		b, err := fetchFileContent(target.Body.File)
		if err != nil {
			return nil, err
		}
		body = b
	}

	if target.Body.Value != "" {
		body = []byte(target.Body.Value)
	}

	cache.Store(cacheKey, body)

	return body, nil
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
