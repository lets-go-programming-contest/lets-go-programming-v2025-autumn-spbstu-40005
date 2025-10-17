package config

import (
	"fmt"
	"strconv"
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
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ConvertFloat(float string) (float64, error) {
	dotFloat := strings.Replace(float, ",", ".", 1)

	result, err := strconv.ParseFloat(dotFloat, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert to float: %w", err)
	}

	return result, nil
}
