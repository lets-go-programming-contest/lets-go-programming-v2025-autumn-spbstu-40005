package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

type Currency struct {
	NumCode  int           `json:"num_code" xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value" xml:"Valute"`
}

type CurrencyData struct {
	Values []Currency `json:"value" xml:"Valute"`
}

func (value *CurrencyValue) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := decoder.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("failed to unmarshal currency value: %w", err)
	}

	str = strings.ReplaceAll(strings.TrimSpace(str), ",", ".")
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("failed to parse currency value '%s': %w", str, err)
	}

	*value = CurrencyValue(result)

	return nil
}
