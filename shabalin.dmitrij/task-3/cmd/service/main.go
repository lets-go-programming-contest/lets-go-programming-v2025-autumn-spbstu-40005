package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dmitei/task-3/internal/launcher"
)

func main() {
	configurationPath := flag.String("config", "configs/configuration.yaml", "Path to YAML configuration file")
	flag.Parse()

	if *configurationPath == "" {
		log.Fatal("configuration path cannot be empty")
	}

	applicationError := launcher.StartApplication(*configurationPath)
	if applicationError != nil {
		fmt.Printf("Application error: %v\n", applicationError)
		panic(applicationError)
	}
}
