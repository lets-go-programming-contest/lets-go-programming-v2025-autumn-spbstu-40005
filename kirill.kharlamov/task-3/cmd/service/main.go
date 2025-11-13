package main

import (
	"flag"
	"sort"

	"kirill.kharlamov/task-3/internal/config"
	"kirill.kharlamov/task-3/internal/currency"
	"kirill.kharlamov/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := parser.LoadYAMLConfig[config.Config](*configPath)
	if err != nil {
		panic(err)
	}

	currencyData, err := parser.LoadXMLData[currency.ValCurs](cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(currencyData.Valute, func(i, j int) bool {
		return currencyData.Valute[i].Value > currencyData.Valute[j].Value
	})

	if err = parser.SaveAsJSON(currencyData.Valute, cfg.OutputFile, cfg.DirPermissions); err != nil {
		panic(err)
	}
}
