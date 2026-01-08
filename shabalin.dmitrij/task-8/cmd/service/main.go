package main

import (
	"fmt"
	"log"

	"github.com/dmiteo/task-8/pkg/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s %s", cfg.Environment, cfg.LogLevel, cfg.AppName, cfg.Version)
}
