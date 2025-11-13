package bank

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyRates struct {
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

type rate float64

type Currency struct {
	NumberCode int    `json:"num_code"  xml:"NumCode"`
	Code       string `json:"char_code" xml:"CharCode"`
	Rate       rate   `json:"value"     xml:"Value"`
}

func (r *rate) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return fmt.Errorf("failed to decode rate value: %w", err)
	}

	parsed, err := strconv.ParseFloat(strings.ReplaceAll(content, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("rate conversion error: %w", err)
	}

	*r = rate(parsed)

	return nil
}
