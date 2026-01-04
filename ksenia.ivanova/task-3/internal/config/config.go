package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigFileNotExist     = errors.New("config file does not exist")
	ErrInputFileNotSpecified  = errors.New("input-file is not specified in config")
	ErrOutputFileNotSpecified = errors.New("output-file is not specified in config")
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (*AppConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotExist, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %w", path, err)
	}

	if config.InputFile == "" {
		return nil, ErrInputFileNotSpecified
	}

	if config.OutputFile == "" {
		return nil, ErrOutputFileNotSpecified
	}

	return &config, nil
}
