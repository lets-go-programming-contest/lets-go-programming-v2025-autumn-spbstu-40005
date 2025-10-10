package main

import (
	"flag"

	"github.com/DariaKhokhryakova/task-3/internal/config"
	"github.com/DariaKhokhryakova/task-3/internal/currency"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the configuration file")
	flag.Parse()

	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := currency.ParseXMLData(config.InputFile)
	if err != nil {
		panic(err)
	}
	result, err := currency.ProcessCurrencies(valCurs)
	if err != nil {
		panic(err)
	}
	err = currency.SaveResults(result, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
