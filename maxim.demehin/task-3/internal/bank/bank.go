package bank

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
	NumCode  int         `xml:"NumCode" json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Value    CustomFloat `xml:"Value" json:"value"`
}

type CustomFloat float64

func (f *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp string

	err := d.DecodeElement(&temp, &start)
	if err != nil {
		return err
	}

	if temp == "" {
		*f = 0.0

		return nil
	}

	valueStr := strings.ReplaceAll(temp, ",", ".")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Value: %w", err)
	}

	*f = CustomFloat(value)

	return nil
}

type ByValue []Valute

func (a ByValue) Len() int {
	return len(a)
}

func (a ByValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByValue) Less(i, j int) bool {
	return a[i].Value > a[j].Value
}

func (v *ValCurs) SortValutes() {
	sort.Sort(ByValue(v.Valutes))
}
