package main

import (
	"flag"
	"os"
	"sort"

	"nikita.brevnov/task-3/internal/bank"
	"nikita.brevnov/task-3/internal/config"
	"nikita.brevnov/task-3/internal/parser"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := parser.LoadYAMLConfig[config.Config](*cfgPath)
	if err != nil {
		panic(err)
	}

	rates, err := parser.LoadXMLData[bank.CurrencyRates](cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(rates.Currencies, func(i, j int) bool {
		return rates.Currencies[i].Rate > rates.Currencies[j].Rate
	})

	if err = parser.SaveAsJSON(rates.Currencies, cfg.OutputFile, os.ModePerm); err != nil {
		panic(err)
	}
}
