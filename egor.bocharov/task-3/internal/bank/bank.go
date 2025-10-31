package bank

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type value float64

type Valute struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    value  `json:"value"     xml:"Value"`
}

func (val *value) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var field string
	if err := decoder.DecodeElement(&field, &start); err != nil {
		return fmt.Errorf("cannot decode value: %w", err)
	}

	result, err := strconv.ParseFloat(strings.ReplaceAll(field, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("convert to float: %w", err)
	}

	*val = value(result)

	return nil
}
