package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func main() {
	configPath := flag.String("config", "", "YAML file required")
	flag.Parse()

	fmt.Println(*configPath)

	cfg, err := parseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg.InputFile, cfg.OutputFile)

}

func parseConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
