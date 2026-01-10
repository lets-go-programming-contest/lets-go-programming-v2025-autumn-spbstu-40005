package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigNotFound = errors.New("config not found")
	ErrConfigRead     = errors.New("config read error")
	ErrConfigInvalid  = errors.New("invalid config format")
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (*AppConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrConfigNotFound, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigRead, err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%w: file %s: %w", ErrConfigInvalid, path, err)
	}

	return &config, nil
}
