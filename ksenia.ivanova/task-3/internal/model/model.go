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
		return fmt.Errorf("failed to parse currency value '%s': %w", str, err)
	}

	*v = CurrencyValue(value)

	return nil
}

func (v CurrencyValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(v), 'f', -1, 64)), nil
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string        `xml:"ID,attr"`
	NumCode  int           `xml:"NumCode"`
	CharCode string        `xml:"CharCode"`
	Nominal  int           `xml:"Nominal"`
	Name     string        `xml:"Name"`
	Value    CurrencyValue `xml:"Value"`
}

type OutputCurrency struct {
	NumCode  int     `json:"iso_num_code"`
	CharCode string  `json:"iso_char_code"`
	Value    float64 `json:"value"`
}

type ByValueDesc []Valute

func (v ByValueDesc) Len() int {
	return len(v)
}

func (v ByValueDesc) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByValueDesc) Less(i, j int) bool {
	return float64(v[i].Value) > float64(v[j].Value)
}
