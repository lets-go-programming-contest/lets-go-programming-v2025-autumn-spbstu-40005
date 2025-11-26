package main

import (
	"flag"

	"github.com/DariaKhokhryakova/task-3/internal/config"
	"github.com/DariaKhokhryakova/task-3/internal/models"
	"github.com/DariaKhokhryakova/task-3/internal/parser"
)

const defaultDirPerm = 0o755

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "path to the configuration file")
	flag.Parse()

	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := parser.ParseXMLData[models.ValCurs](config.InputFile)
	if err != nil {
		panic(err)
	}

	result, err := parser.ProcessCurrencies(valCurs)
	if err != nil {
		panic(err)
	}

	err = parser.SaveJSONResults(result, config.OutputFile, defaultDirPerm)
	if err != nil {
		panic(err)
	}
}
