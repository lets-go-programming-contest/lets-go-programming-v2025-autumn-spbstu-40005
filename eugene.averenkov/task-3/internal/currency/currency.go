package currency

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Valute struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    currencyValue `json:"value"     xml:"Value"`
}

type currencyValue float64

func (v *currencyValue) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := decoder.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}

	str = strings.ReplaceAll(str, ",", ".")

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("failed to parse value: %w", err)
	}

	*v = currencyValue(value)

	return nil
}

func SortByValueDesc(valutes []Valute) {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})
}
