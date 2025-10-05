package main

import (
	"flag"
	"fmt"

	"github.com.P3rCh1/task-3/internal/bank"
	"github.com.P3rCh1/task-3/internal/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "select path to config file")
	flag.Parse()

	config, err := config.ParseFile(*configPath)
	if err != nil {
		panic(fmt.Sprintf("parse config: %s", err))
	}

	bank, err := bank.ParseFile(config.Input)
	if err != nil {
		panic(fmt.Sprintf("parse input-file: %s", err))
	}

	if err := bank.EncodeJSONToFIle(config.Output); err != nil {
		panic(fmt.Sprintf("encode bank: %s", err))
	}
}
