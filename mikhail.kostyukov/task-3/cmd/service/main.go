package main;

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	InputFile string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name	 string `xml:"Name"`
	Value    string `xml:"Value"`
}

const inputFilePath = "1.xml"

func loadConfig(path string) *AppConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Error reading file %s: %w\n", path, err))
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(fmt.Errorf("Error unmarshaling YAML: %w\n", err))
	}

	return &config
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config := loadConfig(*configPath)

	file, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Errorf("Error reading file %s: %w\n", config.InputFile, err))
	}

	//TODO Check error Close

	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var rates ValCurs
	if err := decoder.Decode(&rates); err != nil {
		panic(fmt.Errorf("Error unmarshaling XML: %w\n", err))
	}

	fmt.Printf("Successfully read %d valutes from %s\n", len(rates.Valutes), config.InputFile)
	fmt.Printf("Wrote to %s\n", config.OutputFile)
}
