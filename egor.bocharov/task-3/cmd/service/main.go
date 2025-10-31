package main

import (
	"flag"
	"sort"

	"egor.bocharov/task-3/internal/bank"
	"egor.bocharov/task-3/internal/config"
	"egor.bocharov/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := parser.ParseYAMLFile[config.Config](*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := parser.ParseXMLFile[bank.ValCurs](config.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(valCurs.Valute, func(i, j int) bool {
		return valCurs.Valute[i].Value > valCurs.Valute[j].Value
	})

	if err = parser.EncodeFile(valCurs.Valute, config.OutputFile); err != nil {
		panic(err)
	}
}
