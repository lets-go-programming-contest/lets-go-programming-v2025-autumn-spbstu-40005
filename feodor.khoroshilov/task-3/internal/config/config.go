package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrInputRequired  = errors.New("input-file is required")
	ErrOutputRequired = errors.New("output-file is required")
)

type Settings struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadSettings(configPath string) (*Settings, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	if settings.InputFile == "" {
		return nil, ErrInputRequired
	}

	if settings.OutputFile == "" {
		return nil, ErrOutputRequired
	}

	return &settings, nil
}