package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

//TODO Раскидать все по файлам

//TODO Сократить число структур до минимума

//TODO Исправить капитализацию ошибок (что бы это ни значило)

//TODO Check error defer Close

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
}

type OutputValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type ByValueDesc []Valute

func (v ByValueDesc) Len() int {
	return len(v)
}

func (v ByValueDesc) Swap(i, j int) {
	//TODO add checks
	v[i], v[j] = v[j], v[i]
}

func (v ByValueDesc) Less(i, j int) bool {
	//TODO add checks
	val1 := parseFloatValue(v[i].Value)
	val2 := parseFloatValue(v[j].Value)

	return val1 > val2
}

func parseFloatValue(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	val, err := strconv.ParseFloat(s, 64)

	if err != nil {
		panic(fmt.Errorf("Error parsing value %s: %w", s, err))
	}

	return val
}

func loadConfig(path string) *AppConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Error reading file %s: %w", path, err))
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(fmt.Errorf("Error unmarshaling YAML: %w", err))
	}

	return &config
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config := loadConfig(*configPath)

	file, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Errorf("Error reading file %s: %w", config.InputFile, err))
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var rates ValCurs
	if err := decoder.Decode(&rates); err != nil {
		panic(fmt.Errorf("Error unmarshaling XML: %w", err))
	}

	sort.Sort(ByValueDesc(rates.Valutes))

	outputData := make([]OutputValute, 0, len(rates.Valutes))
	for _, v := range rates.Valutes {
		outputData = append(outputData, OutputValute{v.NumCode, v.CharCode, parseFloatValue(v.Value)})
	}

	jsonData, err := json.MarshalIndent(outputData, "", "\t")
	if err != nil {
		panic(fmt.Errorf("Error marshaling JSON: %w", err))
	}

	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Errorf("Error creating directory %s: %w", outputDir, err))
	}

	if err := os.WriteFile(config.OutputFile, jsonData, 0600); err != nil {
		panic(fmt.Errorf("Error writing file %s: %w", config.OutputFile, err))
	}
}
