package main

import (
	"flag"
	"sort"

	"mohamedamine.drai/task-3/internal/config"
	"mohamedamine.drai/task-3/internal/converter"
	"mohamedamine.drai/task-3/internal/jsonwriter"
	"mohamedamine.drai/task-3/internal/xmlparser"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o600
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	config, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	currencyList := converter.Rates{Data: []converter.Currency{}}

	err = xmlparser.ReadXml(config.InputFile, &currencyList)
	if err != nil {
		panic(err)
	}

	sort.Slice(currencyList.Data, func(i, j int) bool {
		return currencyList.Data[i].Value > currencyList.Data[j].Value
	})

	err = jsonwriter.SaveToJSON(config.OutputFile, currencyList.Data, dirPermissions, filePermissions)
	if err != nil {
		panic(err)
	}
}
