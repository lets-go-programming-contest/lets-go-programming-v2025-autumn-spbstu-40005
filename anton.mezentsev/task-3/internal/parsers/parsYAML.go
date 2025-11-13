package parsers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAML[T any](configfilePath string) (*T, error) {
	file, err := os.Open(configfilePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("closing file: " + err.Error())
		}
	}()

	var config T

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("yaml decoding: %w", err)
	}

	return &config, nil
}
