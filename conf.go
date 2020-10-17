package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Target struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
	Body   string `yaml:"body"`
}

type Config struct {
	Port           string
	ConfigFilePath string

	Targets []Target `yaml:"targets"`
}

func LoadConfig() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	addr := fmt.Sprintf(":%s", port)

	configLocation := os.Getenv("CONFIG")
	if configLocation == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}

		configLocation = path.Join(pwd, "config.yml")
	}

	cfg := Config{Port: addr, ConfigFilePath: configLocation}

	fileBytes, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
