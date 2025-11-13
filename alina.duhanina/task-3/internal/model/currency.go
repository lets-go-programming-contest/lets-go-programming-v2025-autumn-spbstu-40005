package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type CurrencyValue float64

func (cv *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	normalized := strings.ReplaceAll(s, ",", ".")

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*cv = CurrencyValue(value)

	return nil
}

type Valute struct {
	XMLName  xml.Name      `json:"-"         xml:"Valute"`
	ID       string        `json:"-"         xml:"ID,attr"`
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Nominal  int           `json:"-"         xml:"Nominal"`
	Name     string        `json:"-"         xml:"Name"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}
