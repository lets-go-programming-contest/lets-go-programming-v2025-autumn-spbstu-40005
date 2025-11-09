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
		converted, ok := c.convertCurrency(currency)
		if ok {
			output = append(output, converted)
		}
	}

	c.sortByValueDescending(output)
	return output
}

func (c *CurrencyConverter) convertCurrency(currency xmlparser.Currency) (CurrencyOutput, bool) {
	value, err := c.parseValue(currency.Value)
	if err != nil {
		return CurrencyOutput{}, false
	}

	numCode, err := c.parseNumCode(currency.NumCode)
	if err != nil {
		return CurrencyOutput{}, false
	}

	return CurrencyOutput{
		NumCode:  numCode,
		CharCode: strings.TrimSpace(currency.CharCode),
		Value:    value,
	}, true
}

func (c *CurrencyConverter) parseValue(valueStr string) (float64, error) {
	normalized := strings.ReplaceAll(valueStr, ",", ".")
	return strconv.ParseFloat(normalized, 64)
}

func (c *CurrencyConverter) parseNumCode(numCodeStr string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(numCodeStr))
}

func (c *CurrencyConverter) sortByValueDescending(currencies []CurrencyOutput) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}
