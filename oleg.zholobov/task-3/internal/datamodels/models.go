package datamodels

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

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

func SortByValueDesc(valutes []Valute) []Valute {
	sorted := make([]Valute, len(valutes))
	copy(sorted, valutes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}
