package parsers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAML[T any](configfilePath string) (*T, error) {
	file, err := os.Open(configfilePath)
	if err != nil {
		return nil, fmt.Errorf("Error: opening a file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("Error: closing config file")
		}
	}()

	var config T

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Error: yaml decoding: %w", err)
	}

	return &config, nil
}
