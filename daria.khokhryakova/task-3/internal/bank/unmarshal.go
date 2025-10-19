package bank

import "encoding/xml"

func (c *CurrencyResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type xmlValute struct {
		NumCode  string `xlm:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var valute xmlValute

	if err := d.DecodeElement(&valute, &start); err != nil {
		return err
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
