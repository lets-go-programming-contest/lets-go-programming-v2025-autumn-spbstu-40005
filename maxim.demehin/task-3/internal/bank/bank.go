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
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Value    float64 `xml:"Value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type valuteTemp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var temp valuteTemp
	if err := d.DecodeElement(&temp, &start); err != nil {
		return err
	}

	v.CharCode = temp.CharCode

	numCode, err := strconv.Atoi(temp.NumCode)
	if err != nil {
		return fmt.Errorf("failed to parse NumCode '%s': %w", temp.NumCode, err)
	}

	v.NumCode = numCode

	valueStr := strings.ReplaceAll(temp.Value, ",", ".")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Value '%s': %w", temp.Value, err)
	}

	v.Value = value

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

func (v *ValCurs) sortValutes() {
	sort.Sort(ByValue(v.Valutes))
}
