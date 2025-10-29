package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (v *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("failed to unmarshal currency value: %w", err)
	}

	str = strings.ReplaceAll(str, ",", ".")

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("failed to parse currency value: %w", err)
	}

	*v = CurrencyValue(value)

	return nil
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type ByValueDesc []Valute

func (v ByValueDesc) Len() int {
	return len(v)
}

func (v ByValueDesc) Swap(first, second int) {
	if (first < 0) || (first >= len(v)) {
		panic("first index out of range")
	} else if (second < 0) || (second >= len(v)) {
		panic("second index out of range")
	}

	v[first], v[second] = v[second], v[first]
}

func (v ByValueDesc) Less(first, second int) bool {
	if (first < 0) || (first >= len(v)) {
		panic("first index out of range")
	} else if (second < 0) || (second >= len(v)) {
		panic("second index out of range")
	}

	return v[first].Value > v[second].Value
}
