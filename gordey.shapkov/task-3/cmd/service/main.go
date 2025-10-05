package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
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

func main() {
	configPath := flag.String("config", "", "YAML file required")
	flag.Parse()

	fmt.Println(*configPath)

	cfg, err := parseConfigFile(*configPath)
	if err != nil {
		fmt.Println(err)
	}

	valCurs, err := parseXmlFile(cfg.InputFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(valCurs.Date, valCurs.Name)
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
