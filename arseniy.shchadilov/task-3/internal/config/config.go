package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

var (
	ErrInputFileRequired  = errors.New("input-file is required")
	ErrOutputFileRequired = errors.New("output-file is required")
)

func (c *Config) Validate() error {
	if c.InputFile == "" {
		return ErrInputFileRequired
	}

	if c.OutputFile == "" {
		return ErrOutputFileRequired
	}

	return nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
