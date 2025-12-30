package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var errConfigUnmarshal = errors.New("failed to unmarshal config yaml")

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", errConfigUnmarshal, err)
	}

	return &cfg, nil
}
