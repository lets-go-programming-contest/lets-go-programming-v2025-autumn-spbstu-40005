package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TvoyBatyA12343/task-3/internal/bank"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func loadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil

}

func main() {
	cfgPath := flag.String("config", "", "path to config file")
	flag.Parse()

	config, err := loadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := bank.ParseXML(config.InputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	valutes, err := bank.ParseToJSON(valCurs.Valutes)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bank.SaveValutesToFile(valutes, config.OutputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

}
