package main

import (
	"flag"

	"mohamedamine.drai/task-3/internal/config"
	"mohamedamine.drai/task-3/internal/converter"
	"mohamedamine.drai/task-3/internal/jsonwriter"
	"mohamedamine.drai/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("config path not provided")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	data, err := xmlparser.ReadXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	converter := converter.NewCurrencyConverter()
	sortedCurrencies := converter.ConvertAndSort(data.Currencies)

	if err := jsonwriter.SaveToJSON(sortedCurrencies, cfg.OutputFile); err != nil {
		panic(err)
	}
}
