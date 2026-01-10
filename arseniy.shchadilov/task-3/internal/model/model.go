package model

import (
	"encoding/xml"
	"fmt"
	"sort"
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
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type byValueDesc []Valute

func (v byValueDesc) Len() int           { return len(v) }
func (v byValueDesc) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byValueDesc) Less(i, j int) bool { return v[i].Value > v[j].Value }

func (vc *ValCurs) SortByValueDesc() {
	sort.Sort(byValueDesc(vc.Valutes))
}
