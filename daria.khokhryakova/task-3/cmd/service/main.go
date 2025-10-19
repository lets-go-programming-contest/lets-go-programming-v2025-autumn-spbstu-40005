package main

import (
	"flag"

	"github.com/DariaKhokhryakova/task-3/internal/bank"
	"github.com/DariaKhokhryakova/task-3/internal/config"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "path to the configuration file")
	flag.Parse()

	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := bank.ParseXMLData(config.InputFile)
	if err != nil {
		panic(err)
	}

	result, err := bank.ProcessCurrencies(valCurs)
	if err != nil {
		panic(err)
	}

	err = bank.SaveResults(result, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
