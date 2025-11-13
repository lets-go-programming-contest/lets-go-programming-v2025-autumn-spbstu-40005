package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYAMLConfig[T any](configPath string) (*T, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("failed to close config file")
		}
	}()

	var configuration T

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&configuration); err != nil {
		return nil, fmt.Errorf("failed to decode YAML configuration: %w", err)
	}

	return &configuration, nil
}
