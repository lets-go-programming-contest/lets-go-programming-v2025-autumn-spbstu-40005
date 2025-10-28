package main

import (
	"flag"
	"fmt"
	"log"

	"alina.duhanina/task-3/internal/config"
	"alina.duhanina/task-3/internal/converter"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("Config path is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Panicf("Error loading config: %v", err)
	}

	err = converter.ConvertXMLToJSON(cfg.InputFile, cfg.OutputFile)
	if err != nil {
		log.Panicf("Error converting XML to JSON: %v", err)
	}

	fmt.Printf("Program completed successfully. Output file: %s\n", cfg.OutputFile)
}
