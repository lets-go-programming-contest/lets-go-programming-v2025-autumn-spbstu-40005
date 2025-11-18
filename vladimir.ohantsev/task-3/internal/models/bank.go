package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Bank struct {
	Currencies []Currency `json:"valute" xml:"Valute"`
}

type valueType float64

func (v *valueType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var stringField string
	if err := decoder.DecodeElement(&stringField, &start); err != nil {
		return fmt.Errorf("decode value string: %w", err)
	}

	got, err := strconv.ParseFloat(strings.Replace(stringField, ",", ".", 1), 64)
	if err != nil {
		return fmt.Errorf("invalid value type: %w", err)
	}

	*v = valueType(got)

	return nil
}

type Currency struct {
	NumCode  int       `json:"num_code"  xml:"NumCode"`
	CharCode string    `json:"char_code" xml:"CharCode"`
	Value    valueType `json:"value"     xml:"Value"`
}
