package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	Date    string     `xml:"Date,attr"`
	Name    string     `xml:"name,attr"`
	Valutes []Currency `xml:"Valute"`
}

type value float64

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    value  `json:"value"     xml:"Value"`
}

func (val *value) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := decoder.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("cannot decode value: %w", err)
	}

	floatVal, err := ConvertFloat(v)
	if err != nil {
		return fmt.Errorf("convert to float: %w", err)
	}

	*val = value(floatVal)

	return nil
}

func ConvertFloat(float string) (float64, error) {
	dotFloat := strings.Replace(float, ",", ".", 1)

	result, err := strconv.ParseFloat(dotFloat, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert to float: %w", err)
	}

	return result, nil
}
