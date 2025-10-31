package main

import (
	"flag"
	"fmt"
	"log"

	"alina.duhanina/task-3/internal/config"
	"alina.duhanina/task-3/internal/converter"
	"alina.duhanina/task-3/internal/model"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Panicf("Error loading config: %v", err)
	}

	valCurs, err := converter.ParseXML[model.ValCurs](cfg.InputFile)
	if err != nil {
		log.Panicf("Error parsing XML: %v", err)
	}

	currencies := converter.ConvertAndSortCurrencies(valCurs)

	if err := converter.SaveJSON(cfg.OutputFile, currencies); err != nil {
		log.Panicf("Error saving JSON: %v", err)
	}

	err = converter.ConvertXMLToJSON(cfg.InputFile, cfg.OutputFile)
	if err != nil {
		log.Panicf("Error converting XML to JSON: %v", err)
	}

	fmt.Printf("Program completed successfully. Output file: %s\n", cfg.OutputFile)
}
