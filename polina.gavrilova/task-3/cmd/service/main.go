package main

import (
	"flag"

	"polina.gavrilova/task-3/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("--config flag is required")
	}

	_, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}
}