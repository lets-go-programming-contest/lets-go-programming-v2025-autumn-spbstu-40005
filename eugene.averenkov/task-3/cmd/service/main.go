package main

import (
	"flag"
	"log"

	"currency-converter/internal/config"
	"currency-converter/internal/converter"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	valutes, err := converter.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	converter.SortByValueDesc(valutes)

	if err := converter.WriteToJSON(valutes, cfg.OutputFile); err != nil {
		panic(err)
	}

	log.Printf("Successfully processed %d currencies. Output: %s", len(valutes), cfg.OutputFile)
}
