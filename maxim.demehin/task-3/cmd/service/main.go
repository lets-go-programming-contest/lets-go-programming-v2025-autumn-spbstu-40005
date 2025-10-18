package main

import (
	"flag"

	"github.com/TvoyBatyA12343/task-3/internal/jsonutils"
	"github.com/TvoyBatyA12343/task-3/internal/parser"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file")
	flag.Parse()

	config, err := parser.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	valutes, err := parser.ParseXML(config.InputFile)
	if err != nil {
		panic(err)
	}

	err = jsonutils.SaveValutesToFile(valutes, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
