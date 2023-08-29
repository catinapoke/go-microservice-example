package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	ServicePort string `yaml:"port"`
	PostgresUrl string `yaml:"postgresURL"`
	Redis       struct {
		Url string `yaml:"url"`
		Key string `yaml:"key"`
	} `yaml:"redis"`
	Nats struct {
		Url     string `yaml:"url"`
		Subject string `yaml:"subject"`
	} `yaml:"nats"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}
