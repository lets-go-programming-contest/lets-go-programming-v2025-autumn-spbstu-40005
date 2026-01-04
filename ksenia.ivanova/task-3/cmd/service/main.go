package main

import (
	"flag"
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
		panic(err)
	}

	var rates model.ValCurs
	if err := converter.ParseXMLFile(appConfig.InputFile, &rates); err != nil {
		panic(err)
	}

	sort.Sort(model.ByValueDesc(rates.Valutes))

	outputCurrencies := make([]model.OutputCurrency, len(rates.Valutes))
	for i, valute := range rates.Valutes {
		outputCurrencies[i] = model.OutputCurrency{
			NumCode:  int(valute.NumCode),
			CharCode: valute.CharCode,
			Value:    float64(valute.Value),
		}
	}

	if err := converter.WriteToJSON(outputCurrencies, appConfig.OutputFile); err != nil {
		panic(err)
	}
}
