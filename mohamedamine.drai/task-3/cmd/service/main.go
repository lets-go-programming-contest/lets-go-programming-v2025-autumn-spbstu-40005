package main

import (
	"flag"

	"mohamedamine.drai/task-3/internal/config"
	"mohamedamine.drai/task-3/internal/converter"
	"mohamedamine.drai/task-3/internal/jsonwriter"
	"mohamedamine.drai/task-3/internal/xmlparser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	xmlData, err := xmlparser.ReadXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	conv := converter.NewCurrencyConverter()
	out := conv.ConvertAndSort(xmlData.Valutes)

	if err := jsonwriter.SaveToJSON(out, cfg.OutputFile); err != nil {
		panic(err)
	}
}
