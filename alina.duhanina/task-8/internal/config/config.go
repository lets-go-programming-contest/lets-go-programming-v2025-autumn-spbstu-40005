package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

//go:embed dev.yaml
var devConfigData []byte

//go:embed prod.yaml
var prodConfigData []byte

func loadConfig(data []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetConfig() (*Config, error) {
	return loadConfig(prodConfigData)
}
