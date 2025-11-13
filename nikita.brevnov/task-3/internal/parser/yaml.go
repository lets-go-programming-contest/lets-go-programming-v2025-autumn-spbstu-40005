package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYAMLConfig[T any](filePath string) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var settings T
	if err = yaml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("YAML parsing failed: %w", err)
	}

	return &settings, nil
}
