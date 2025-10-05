package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"path/filepath"
	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"sort"
	"strings"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    string `json:"value"`
}

func main() {
	configPath := flag.String("config", "", "YAML file required")
	flag.Parse()

	cfg, err := parseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
	}

	valCurs, err := parseXmlFile(cfg.InputFile)
	if err != nil {
		fmt.Println(err)
	}

	valutes := valCurs.Valutes
	currencies := convertToJson(valutes)
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
	err = saveToJSON(currencies, cfg.OutputFile)
	if err != nil {
                fmt.Println(err)
        }
}

func parseConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseXmlFile(path string) (*ValCurs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = createCharsetReader

	valCurs := &ValCurs{}
	err = decoder.Decode(valCurs)
	if err != nil {
		return nil, err
	}

	return valCurs, nil
}

func createCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1251" {
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}
	return input, nil
}

func convertFloat(float string) string {
	partsOfFloat := strings.Split(float, ",")
	return partsOfFloat[0] + "." + partsOfFloat[1]
}

func convertToJson(valutes []Valute) []Currency {
	currencies := make([]Currency, len(valutes))

	for idx, valute := range valutes {
		value := convertFloat(valute.Value)
		currencies[idx] = Currency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
	}

	return currencies
}

func saveToJSON(currencies []Currency, outputPath string) error {
	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}
