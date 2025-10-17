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
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type valuteTemp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var temp valuteTemp

	err := d.DecodeElement(&temp, &start)
	if err != nil {
		return err
	}

	v.CharCode = temp.CharCode

	numCode, err := strconv.Atoi(temp.NumCode)
	if err != nil {
		return fmt.Errorf("failed to parse NumCode: %w", err)
	}

	v.NumCode = numCode

	valueStr := strings.ReplaceAll(temp.Value, ",", ".")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Value: %w", err)
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
