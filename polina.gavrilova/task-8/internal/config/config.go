package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	err := yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	return &cfg, nil
}
