package currency

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

func (v *value) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	parsed, err := strconv.ParseFloat(strings.ReplaceAll(content, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("failed to parse float value: %w", err)
	}

	*v = value(parsed)

	return nil
}
