package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	SourceFile      string `yaml:"input-file"`
	DestinationFile string `yaml:"output-file"`
}

func LoadApplicationConfig(configFilePath string) (*ApplicationConfig, error) {
	configurationContent, readError := os.ReadFile(configFilePath)
	if readError != nil {
		return nil, fmt.Errorf("cannot read configuration file at %q: %w", configFilePath, readError)
	}

	appConfig := new(ApplicationConfig)

	if unmarshalError := yaml.Unmarshal(configurationContent, appConfig); unmarshalError != nil {
		return nil, fmt.Errorf("cannot unmarshal YAML configuration: %w", unmarshalError)
	}

	return appConfig, nil
}
