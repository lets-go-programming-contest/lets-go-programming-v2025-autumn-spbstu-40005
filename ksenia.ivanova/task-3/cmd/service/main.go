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

	filtered := make([]model.Valute, 0, len(rates.Valutes))
	for _, v := range rates.Valutes {
		if v.NumCode > 0 {
			filtered = append(filtered, v)
		}
	}

	sort.Sort(model.ByNumCode(filtered))

	output := make([]model.OutputCurrency, len(filtered))
	for i, v := range filtered {
		output[i] = model.OutputCurrency{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    float64(v.Value),
		}
	}

	if err := converter.WriteToJSON(output, appConfig.OutputFile); err != nil {
		log.Fatal(err)
	}
}
