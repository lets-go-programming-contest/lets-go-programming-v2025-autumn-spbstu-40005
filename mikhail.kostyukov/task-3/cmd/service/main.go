package main;

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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
	//TODO delete Nominal & Name
	NumCode  int    `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Nominal  int    `xml:"Nominal" json:"-"`
	Name	 string `xml:"Name" json:"-"`
	Value    string `xml:"Value" json:"value"`
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
		panic(fmt.Errorf("Error parsing value %s: %w\n", s, err))
	}

	return val
}

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

	defer file.Close() 	//TODO Check error defer Close

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var rates ValCurs
	if err := decoder.Decode(&rates); err != nil {
		panic(fmt.Errorf("Error unmarshaling XML: %w\n", err))
	}

	sort.Sort(ByValueDesc(rates.Valutes))

	fmt.Printf("Successfully read %d valutes from %s\n", len(rates.Valutes), config.InputFile)

	for _, v := range rates.Valutes {
		fmt.Printf("%d %s %s\n", v.NumCode, v.CharCode, v.Value)
	}
}
