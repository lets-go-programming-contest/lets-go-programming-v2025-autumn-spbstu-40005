package main

import (
	"flag"

	"arseniy.shchadilov/task-3/internal/config"
	"arseniy.shchadilov/task-3/internal/converter"
	"arseniy.shchadilov/task-3/internal/model"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("--config flag is required")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	var currencyData model.ValCurs
	if err := converter.ParseXMLFile(cfg.InputFile, &currencyData); err != nil {
		panic(err)
	}

	currencyData.SortByValueDesc()

	if err := converter.WriteToJSON(currencyData.Valutes, cfg.OutputFile); err != nil {
		panic(err)
	}
}
