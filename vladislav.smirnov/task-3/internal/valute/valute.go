package valute

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCursXML struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    Decimal `json:"value"     xml:"Value"`
}

type Decimal float64

func (dec *Decimal) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string

	if err := decoder.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("fail to decode element: %w", err)
	}

	str = strings.Replace(str, ",", ".", 1)
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return fmt.Errorf("fail to parse float: %w", err)
	}

	*dec = Decimal(value)

	return nil
}


