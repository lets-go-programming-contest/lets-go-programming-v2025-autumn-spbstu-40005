package main

import (
	"encoding/xml"
	"flag"
	"os"

	"eugene.averenkov/task-3/internal/config"
	"eugene.averenkov/task-3/internal/converter"
	"eugene.averenkov/task-3/internal/currency"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
	dirPermissions    = os.ModePerm
	filePermissions   = 0o644
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := converter.ParseXMLFile[struct {
		XMLName xml.Name          `xml:"ValCurs"`
		Valutes []currency.Valute `xml:"Valute"`
	}](cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currency.SortByValueDesc(valCurs.Valutes)

	if err := converter.WriteToJSON(valCurs.Valutes, cfg.OutputFile, dirPermissions, filePermissions); err != nil {
		panic(err)
	}
}
