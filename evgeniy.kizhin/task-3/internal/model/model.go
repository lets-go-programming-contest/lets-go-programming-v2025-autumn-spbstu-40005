package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Date    string     `xml:"Date,attr"`
	Name    string     `xml:"name,attr"`
	Valute  []Currency `xml:"Valute"`
}

type Amount float64

func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	str = strings.ReplaceAll(str, ",", ".")

	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("parse float '%s': %w", str, err)
	}

	*a = Amount(v)

	return nil
}

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    Amount `json:"value"     xml:"Value"`
}
