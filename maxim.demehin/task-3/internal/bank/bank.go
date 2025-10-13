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
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type ValCursRaw struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valutes []ValuteRaw `xml:"Valute"`
}

type ValuteRaw struct {
	ID        string `xml:"ID,attr"`
	NumCode   int    `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

type Valute struct {
	ID        string  `xml:"ID,attr"`
	NumCode   int     `xml:"NumCode"`
	CharCode  string  `xml:"CharCode"`
	Nominal   int     `xml:"Nominal"`
	Name      string  `xml:"Name"`
	Value     float64 `xml:"Value"`
	VunitRate string  `xml:"VunitRate"`
}

type ValuteToOut struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func (v *Valute) convertFromRaw(valuteRaw ValuteRaw) error {
	v.ID = valuteRaw.ID
	v.NumCode = valuteRaw.NumCode
	v.CharCode = valuteRaw.CharCode
	v.Nominal = valuteRaw.Nominal
	v.Name = valuteRaw.Name
	v.VunitRate = valuteRaw.VunitRate

	valueStr := strings.Replace(valuteRaw.Value, ",", ".", -1)

	floatVal, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed conversion from ValuteRaw to Valute: %w", err)
	}

	v.Value = floatVal

	return nil
}

func (v *ValCurs) convertFromRaw(valCursRaw ValCursRaw) error {
	temp := ValCurs{
		XMLName: valCursRaw.XMLName,
		Date:    valCursRaw.Date,
		Name:    valCursRaw.Name,
	}

	temp.Valutes = make([]Valute, len(valCursRaw.Valutes))

	for index := range valCursRaw.Valutes {
		err := temp.Valutes[index].convertFromRaw(valCursRaw.Valutes[index])
		if err != nil {
			return fmt.Errorf("failed conversion from ValCursRaw to ValCurs: %w", err)
		}
	}

	*v = temp
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
