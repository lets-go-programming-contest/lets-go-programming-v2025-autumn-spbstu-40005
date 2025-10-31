package parsers

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrOpeningFile     = errors.New("error opening file")
	ErrClosingFileYAML = errors.New("error closing file")
	ErrYAMLDecoding    = errors.New("error yaml decoding")
)

func ParseYAML[T any](configfilePath string) (*T, error) {
	file, err := os.Open(configfilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpeningFile, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(ErrClosingFileYAML)
		}
	}()

	var config T

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrYAMLDecoding, err)
	}

	return &config, nil
}
