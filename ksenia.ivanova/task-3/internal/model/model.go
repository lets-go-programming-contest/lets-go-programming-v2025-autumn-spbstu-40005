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

	str = strings.TrimSpace(str)
	if str == "" {
		return fmt.Errorf("currency value is empty")
	}

	str = strings.ReplaceAll(str, ",", ".")
	str = strings.ReplaceAll(str, " ", "")

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
	ID       string        `xml:"-"`
	NumCode  int           `xml:"-"`
	CharCode string        `xml:"CharCode"`
	Nominal  int           `xml:"Nominal"`
	Name     string        `xml:"Name"`
	Value    CurrencyValue `xml:"Value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "ID" {
			v.ID = attr.Value
			break
		}
	}

	type ValuteInner struct {
		NumCode  string        `xml:"NumCode"`
		CharCode string        `xml:"CharCode"`
		Nominal  int           `xml:"Nominal"`
		Name     string        `xml:"Name"`
		Value    CurrencyValue `xml:"Value"`
	}

	var inner ValuteInner
	if err := d.DecodeElement(&inner, &start); err != nil {
		return err
	}

	v.CharCode = inner.CharCode
	v.Nominal = inner.Nominal
	v.Name = inner.Name
	v.Value = inner.Value

	s := strings.TrimSpace(inner.NumCode)
	if s == "" {
		return fmt.Errorf("NumCode is empty")
	}
	if len(s) > 1 {
		s = strings.TrimLeft(s, "0")
		if s == "" {
			s = "0"
		}
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid NumCode '%s': %w", inner.NumCode, err)
	}
	v.NumCode = num

	return nil
}

type OutputCurrency struct {
	NumCode  int     `json:"iso_num_code"`
	CharCode string  `json:"iso_char_code"`
	Value    float64 `json:"value"`
}

type ByNumCode []Valute

func (v ByNumCode) Len() int {
	return len(v)
}

func (v ByNumCode) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByNumCode) Less(i, j int) bool {
	return v[i].NumCode < v[j].NumCode
}

type ByNumCodeDesc []Valute

func (v ByNumCodeDesc) Len() int {
	return len(v)
}

func (v ByNumCodeDesc) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByNumCodeDesc) Less(i, j int) bool {
	return v[i].NumCode > v[j].NumCode
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
