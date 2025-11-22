package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrMissingPaths = errors.New("config file is missing required paths")

type Config struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func (c *Config) Validate() error {
	if c.Input == "" || c.Output == "" {
		return ErrMissingPaths
	}

	return nil
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", filePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config from %q: %w", filePath, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%q: %w", filePath, err)
	}

	return &cfg, nil
}
