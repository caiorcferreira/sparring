package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Body struct {
	Value string `yaml:"value"`
	File  string `yaml:"file"`
}

type Target struct {
	Method       string        `yaml:"method"`
	Path         string        `yaml:"path"`
	Body         Body          `yaml:"body"`
	ResponseTime time.Duration `yaml:"responseTime"`
	StatusCode   int           `yaml:"statusCode"`
}

type Config struct {
	ConfigFilePath string

	Port     string   `yaml:"port"`
	LogLevel string   `yaml:"logLevel"`
	Targets  []Target `yaml:"targets"`
}

func LoadConfig() (Config, error) {
	port := os.Getenv("SPARRING_PORT")
	if port == "" {
		port = "9000"
	}

	configLocation := os.Getenv("SPARRING_CONFIG")
	if configLocation == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}

		configLocation = path.Join(pwd, "config.yml")
	}

	logLevel := os.Getenv("SPARRING_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	cfg := Config{
		Port:           port,
		ConfigFilePath: configLocation,
		LogLevel:       logLevel,
	}

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
