package main;

import (
	"encoding/xml"
	"fmt"
	"os"
	"golang.org/x/net/html/charset"
)

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

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		panic(fmt.Errorf("Error reading file %s: %w\n", inputFilePath, err))
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var rates ValCurs
	if err := decoder.Decode(&rates); err != nil {
		panic(fmt.Errorf("Error unmarshaling XML: %w\n", err))
	}

	fmt.Printf("Successfully unmarshaled %d Valutes.\n", len(rates.Valutes))
	for _, currency := range rates.Valutes {
		fmt.Printf("Code: %s, Nominal: %d, Name: %s, Value: %s\n", currency.CharCode, currency.Nominal, currency.Name, currency.Value)
	}
}
