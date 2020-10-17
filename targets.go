package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
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

		body, err := mountTargetBody(target)
		if err != nil {
			ctx.Logger().Errorf("failed mounting target body for response: %s", err)
			return err
		}

		err = ctx.JSONBlob(http.StatusOK, body)
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
