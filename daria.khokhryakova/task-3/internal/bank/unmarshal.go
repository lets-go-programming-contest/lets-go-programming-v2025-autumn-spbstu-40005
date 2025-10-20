package bank

import (
	"encoding/xml"
	"fmt"
)

func (c *CurrencyResult) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type xmlValute struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var valute xmlValute

	if err := decoder.DecodeElement(&valute, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	numCode, err := convertNumCode(valute.NumCode)
	if err != nil {
		return err
	}

	value, err := convertValue(valute.Value)
	if err != nil {
		return err
	}

	c.NumCode = numCode
	c.CharCode = valute.CharCode
	c.Value = value

	return nil
}
