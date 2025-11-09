package converter

import (
	"sort"
	"strconv"
	"strings"

	"mohamedamine.drai/task-3/internal/xmlparser"
)

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type CurrencyConverter struct{}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{}
}

func (c *CurrencyConverter) ConvertAndSort(currencies []xmlparser.Currency) []CurrencyOutput {
	output := make([]CurrencyOutput, 0, len(currencies))

	for _, currency := range currencies {
		converted := c.convertCurrency(currency)
		output = append(output, converted)
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output
}

func (c *CurrencyConverter) convertCurrency(currency xmlparser.Currency) CurrencyOutput {
	value, err := strconv.ParseFloat(strings.ReplaceAll(currency.Value, ",", "."), 64)
	if err != nil {
		value = 0
	}

	numCode, err := strconv.Atoi(strings.TrimSpace(currency.NumCode))
	if err != nil {
		numCode = 0
	}

	return CurrencyOutput{
		NumCode:  numCode,
		CharCode: strings.TrimSpace(currency.CharCode),
		Value:    value,
	}
}
