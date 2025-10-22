package main

import (
	"flag"
	"sort"

	"github.com/TvoyBatyA12343/task-3/internal/datamodels"
	"github.com/TvoyBatyA12343/task-3/internal/jsonutils"
	"github.com/TvoyBatyA12343/task-3/internal/parserxml"
	"github.com/TvoyBatyA12343/task-3/internal/parseryaml"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := parseryaml.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	valutes, err := parserxml.ParseXML(config.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Sort(datamodels.ByValue(valutes))

	err = jsonutils.SaveValutesToFile(valutes, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
