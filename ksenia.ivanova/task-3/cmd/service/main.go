package main

import (
	"flag"
	"log"
	"sort"

	"ksenia.ivanova/task-3/internal/config"
	"ksenia.ivanova/task-3/internal/converter"
	"ksenia.ivanova/task-3/internal/model"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	appConfig, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var rates model.ValCurs
	if err := converter.ParseXMLFile(appConfig.InputFile, &rates); err != nil {
		log.Fatal(err)
	}

	if len(rates.Valutes) == 0 {
		log.Fatal("input XML contains no Valute entries")
	}

	sort.Sort(model.ByNumCode(rates.Valutes))

	outputCurrencies := make([]model.OutputCurrency, len(rates.Valutes))
	for i, valute := range rates.Valutes {
		outputCurrencies[i] = model.OutputCurrency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    float64(valute.Value),
		}
	}

	if err := converter.WriteToJSON(outputCurrencies, appConfig.OutputFile); err != nil {
		log.Fatal(err)
	}
}
