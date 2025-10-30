package main

import (
	"flag"

	"alexandra.karnauhova/task-1/internal/config"
)

func main() {
	thisConfig := flag.String("config", "non", "Select a configuration file")
	flag.Parse()

	if *thisConfig == "non" {
		panic("Config file is uncorrect")
	}

	cfg, err := config.LoadConfig(*thisConfig)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
}
