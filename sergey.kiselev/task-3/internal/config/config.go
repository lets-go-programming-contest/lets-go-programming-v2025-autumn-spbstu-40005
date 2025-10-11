package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ParseFile(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing config file")
		}
	}()

	return parseYaml(file)
}

func parseYaml(reader io.Reader) (*Config, error) {
	config := new(Config)
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	return config, nil
}
