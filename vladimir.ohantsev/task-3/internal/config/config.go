package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func ParseFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Sprintf("close file: %s", err))
		}
	}()

	return Parse(file)
}

func Parse(r io.Reader) (*Config, error) {
	config := new(Config)
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, fmt.Errorf("decoding config file: %w", err)
	}

	return config, nil
}
