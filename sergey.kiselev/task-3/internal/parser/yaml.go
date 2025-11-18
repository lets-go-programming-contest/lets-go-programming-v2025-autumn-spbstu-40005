package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAMLFile[T any](configPath string) (*T, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing config file")
		}
	}()

	var config T

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	return &config, nil
}
