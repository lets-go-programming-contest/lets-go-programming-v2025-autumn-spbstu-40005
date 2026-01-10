package main

import (
	"flag"
	"sort"

	"ksenia.ivanova/task-3/internal/config"
	"ksenia.ivanova/task-3/internal/jsonwriter"
	"ksenia.ivanova/task-3/internal/xmlparser"
)

const (
	flagConfigName    = "config"
	flagConfigDefault = "config.yaml"
	flagConfigUsage   = "path to config file"
)

func main() {
	configPath := flag.String(flagConfigName, flagConfigDefault, flagConfigUsage)
	flag.Parse()

	appConfig, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	curensies, err := xmlparser.ParseFile(appConfig.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(curensies.Values, func(i, j int) bool {
		return curensies.Values[i].Value > curensies.Values[j].Value
	})

	if err := jsonwriter.Dump(curensies, appConfig.OutputFile); err != nil {
		panic(err)
	}
}
