package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type valuteTemp struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp valuteTemp
	if err := d.DecodeElement(&temp, &start); err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	if temp.NumCode == "" || temp.CharCode == "" || temp.Value == "" {
		return nil
	}

	numCode, err := strconv.Atoi(temp.NumCode)
	if err != nil {
		return fmt.Errorf("invalid NumCode %s: %w", temp.NumCode, err)
	}

	valueStr := strings.Replace(temp.Value, ",", ".", 1)
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("invalid Value %s: %w", temp.Value, err)
	}

	v.NumCode = numCode
	v.CharCode = temp.CharCode
	v.Value = value

	return nil
}
