package main

import (
	"flag"
	"fmt"

	"github.com/deonik3/task-3/internal/bank"
	"github.com/deonik3/task-3/internal/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := config.ParseFile(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := bank.ParseXMLFile(config.InputFile)
	if err != nil {
		panic(err)
	}

	currencies, err := bank.ConvertAndSort(valCurs.Valute)
	if err != nil {
		fmt.Println(err)
	}

	if err = bank.EncodeFile(currencies, config.OutputFile); err != nil {
		panic(err)
	}
}
