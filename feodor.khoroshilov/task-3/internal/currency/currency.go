package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ExchangeData struct {
	XMLName    xml.Name `xml:"ValCurs"`
	DateValue  string   `xml:"Date,attr"`
	SourceName string   `xml:"name,attr"`
	Items      []Item   `xml:"Valute"`
}

type Item struct {
	ID        string  `json:"-"         xml:"ID,attr"`
	NumCode   int     `json:"num_code"  xml:"NumCode"`
	CharCode  string  `json:"char_code" xml:"CharCode"`
	Nominal   int     `json:"-"         xml:"Nominal"`
	ItemName  string  `json:"-"         xml:"Name"`
	RateValue float64 `json:"value"     xml:"Value"`
}

type moneyValue float64

func (mv *moneyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var text string
	if err := d.DecodeElement(&text, &start); err != nil {
		return fmt.Errorf("error decoding value: %w", err)
	}

	if text == "" {
		*mv = moneyValue(0)

		return nil
	}

	text = strings.ReplaceAll(text, ",", ".")

	num, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("error parsing value '%s': %w", text, err)
	}

	*mv = moneyValue(num)

	return nil
}

type itemXML struct {
	ID       string     `xml:"ID,attr"`
	NumCode  string     `xml:"NumCode"`
	CharCode string     `xml:"CharCode"`
	Nominal  int        `xml:"Nominal"`
	Name     string     `xml:"Name"`
	Value    moneyValue `xml:"Value"`
}

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp itemXML
	if err := d.DecodeElement(&temp, &start); err != nil {
		return fmt.Errorf("error decoding XML element: %w", err)
	}

	numCode := 0

	if temp.NumCode != "" {
		var err error
		numCode, err = strconv.Atoi(temp.NumCode)

		if err != nil {
			return fmt.Errorf("failed to parse NumCode '%s': %w", temp.NumCode, err)
		}
	}

	i.ID = temp.ID
	i.NumCode = numCode
	i.CharCode = temp.CharCode
	i.Nominal = temp.Nominal
	i.ItemName = temp.Name
	i.RateValue = float64(temp.Value)

	return nil
}
