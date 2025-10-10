package config

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    InputFile  string `yaml:"input-file"`
    OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*Config, error) {
    _, err := os.Stat(configPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("file does not exist")
        }
    }
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("file cannot be read")
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("did not find expected key")
    }
    if config.InputFile == "" {
        return nil, fmt.Errorf("inputFile field is empty")
    }
    if config.OutputFile == "" {
        return nil, fmt.Errorf("outputFile field is empty")
    }
    return &config, nil
}
