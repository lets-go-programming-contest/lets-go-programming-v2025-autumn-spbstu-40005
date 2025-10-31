package main

import (
	"flag"

	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/processor"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("--config flag is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = processor.Run(cfg)
	if err != nil {
		panic(err)
	}
}
