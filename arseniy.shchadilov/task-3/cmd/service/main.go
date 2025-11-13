package main

import (
	"flag"
	"fmt"

	"arseniy.shchadilov/task-3/internal/config"
	"arseniy.shchadilov/task-3/internal/converter"
	"arseniy.shchadilov/task-3/internal/model"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %v", err))
	}

	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("ERROR: Invalid config: %v", err))
	}

	var currencyData model.ValCurs
	if err := converter.ParseXMLFile(cfg.InputFile, &currencyData); err != nil {
		panic(fmt.Sprintf("ERROR: Failed to parse XML: %v", err))
	}

	currencyData.SortByValueDesc()

	if err := converter.WriteToJSON(currencyData.Valutes, cfg.OutputFile); err != nil {
		panic(fmt.Sprintf("ERROR: Failed to write JSON: %v", err))
	}

	fmt.Printf("Successfully processed %d currency records to %s\n",
		len(currencyData.Valutes), cfg.OutputFile)
}
