package datamodels

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int         `json:"num_code"  xml:"NumCode"`
	CharCode string      `json:"char_code" xml:"CharCode"`
	Value    CustomFloat `json:"value"     xml:"Value"`
}

type CustomFloat float64

func (f *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp string

	err := d.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	valueStr := strings.ReplaceAll(temp, ",", ".")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Value: %w", err)
	}

	*f = CustomFloat(value)

	return nil
}

type ByValue []Valute

func (a ByValue) Len() int {
	return len(a)
}

func (a ByValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByValue) Less(i, j int) bool {
	return a[i].Value > a[j].Value
}
