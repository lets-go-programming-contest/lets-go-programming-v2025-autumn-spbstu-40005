package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/P3rCh1/task-3/internal/jsonwrap"
	"github.com/P3rCh1/task-3/internal/models"
	"github.com/P3rCh1/task-3/internal/xmlwrap"
	"github.com/P3rCh1/task-3/internal/yamlwrap"
)

const permissions = 0o755

func must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", op, err))
	}
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := yamlwrap.ParseFile[models.Config](*configPath)

	must("parse config", err)

	bank, err := xmlwrap.ParseFile[models.Bank](config.Input)

	must("parse input-file", err)

	sort.Slice(
		bank.Currencies,
		func(i, j int) bool {
			return bank.Currencies[i].Value > bank.Currencies[j].Value
		},
	)

	must("encode bank", jsonwrap.EncodeFile(bank.Currencies, config.Output, permissions))
}
