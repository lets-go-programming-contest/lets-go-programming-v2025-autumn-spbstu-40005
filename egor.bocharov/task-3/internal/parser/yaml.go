package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAMLFile[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing file")
		}
	}()

	var result T

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	return &result, nil
}
