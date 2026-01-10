package main

import (
	"flag"
	"fmt"

	"oleg.zholobov/task-3/internal/config"
	"oleg.zholobov/task-3/internal/datamodels"
	"oleg.zholobov/task-3/internal/jsonwriter"
	"oleg.zholobov/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("Config path is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading config: %v", err))
	}

	valutes, err := xmlparser.ParseXML(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing XML: %v", err))
	}

	sortedValutes := datamodels.SortByValueDesc(valutes)

	if err := jsonwriter.SaveJSON(cfg.OutputFile, sortedValutes); err != nil {
		panic(fmt.Sprintf("Error saving JSON: %v", err))
	}

	fmt.Printf("Successfully processed %d currencies. Output saved to: %s\n", len(sortedValutes), cfg.OutputFile)
}
