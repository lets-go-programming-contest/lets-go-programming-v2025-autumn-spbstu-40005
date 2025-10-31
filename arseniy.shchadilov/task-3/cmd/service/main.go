package main

import (
	"flag"
	"fmt"
	"os"

	"arseniy.shchadilov/task-3/internal/config"
	"arseniy.shchadilov/task-3/internal/converter"
	"arseniy.shchadilov/task-3/internal/model"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --config flag is required")
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid config: %v\n", err)
		os.Exit(1)
	}

	var currencyData model.ValCurs
	if err := converter.ParseXMLFile(cfg.InputFile, &currencyData); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse XML: %v\n", err)
		os.Exit(1)
	}

	currencyData.SortByValueDesc()

	if err := converter.WriteToJSON(currencyData.Valutes, cfg.OutputFile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d currency records to %s\n",
		len(currencyData.Valutes), cfg.OutputFile)
}
