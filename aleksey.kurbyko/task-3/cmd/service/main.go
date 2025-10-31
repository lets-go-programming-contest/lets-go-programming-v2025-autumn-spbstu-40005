package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"os"
	"path/filepath"

	"github.com/aleksey.kurbyko/task-3/internal/dataprocessor"
	"github.com/aleksey.kurbyko/task-3/internal/currencyhandler"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML config")
	flag.Parse()

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var paths dataprocessor.FilePaths
	err = yaml.Unmarshal(configData, &paths)
	if err != nil {
		panic(err)
	}

	inputFile, err := os.Open(paths.Input)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	xmlDecoder := xml.NewDecoder(inputFile)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var currencies currencyhandler.CurrencyList
	err = xmlDecoder.Decode(&currencies)
	if err != nil {
		panic(err)
	}

	currencyhandler.SortCurrencies(&currencies)
	jsonOutput := dataprocessor.ConvertToJSON(currencies)

	outputData, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Dir(paths.Output), 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(paths.Output, outputData, 0644)
	if err != nil {
		panic(err)
	}
}
