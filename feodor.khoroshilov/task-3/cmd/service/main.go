package main

import (
	"flag"

	"feodor.khoroshilov/task-3/internal/config"
	"feodor.khoroshilov/task-3/internal/convert"
	"feodor.khoroshilov/task-3/internal/currency"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

const (
	DirPerms = 0o755
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	settings, err := config.LoadSettings(*configPath)
	if err != nil {
		panic(err)
	}

	fulldata, err := convert.LoadXMLData[currency.ExchangeData](settings.InputFile)
	if err != nil {
		panic(err)
	}

	convert.SortItemsByRate(&fulldata.Items)

	if err := convert.SaveItemsAsJSON(fulldata.Items, settings.OutputFile, DirPerms); err != nil {
		panic(err)
	}
}
