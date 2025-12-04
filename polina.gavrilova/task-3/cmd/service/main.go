package main

import (
	"flag"
	"path/filepath"

	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/processor"
)

func main() {
	defaultConfig := filepath.Join("config", "config.yaml")
	configPath := flag.String("config", defaultConfig, "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = processor.Run(cfg)
	if err != nil {
		panic(err)
	}
}
