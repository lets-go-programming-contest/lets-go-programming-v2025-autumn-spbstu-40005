package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config структура для хранения конфигурации
type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func getConfig(dataYaml []byte) (*Config, error) {
	var it Config

	err := yaml.Unmarshal(dataYaml, &it)
	if err != nil {
		return nil, fmt.Errorf("error in parsing config: %v", err)
	}

	return &it, nil
}

func (c *Config) printToConfig() {
	fmt.Printf("%s %s\n", c.Environment, c.LogLevel)
}
