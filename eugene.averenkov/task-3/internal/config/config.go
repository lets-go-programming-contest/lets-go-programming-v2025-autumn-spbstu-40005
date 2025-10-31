package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrInputFileRequired  = errors.New("input-file is required")
	ErrOutputFileRequired = errors.New("output-file is required")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if cfg.InputFile == "" {
		return nil, ErrInputFileRequired
	}

	if cfg.OutputFile == "" {
		return nil, ErrOutputFileRequired
	}

	return &cfg, nil
}
