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
	NumCode   int     `json:"num_code"  xml:"NumCode"`
	CharCode  string  `json:"char_code" xml:"CharCode"`
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

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var temp Item
	if err := d.DecodeElement(&temp, &start); err != nil {
		return fmt.Errorf("error decoding XML element: %w", err)
	}

	numCode := 0

	i.NumCode = numCode
	i.CharCode = temp.CharCode
	i.RateValue = float64(temp.RateValue)

	return nil
}
