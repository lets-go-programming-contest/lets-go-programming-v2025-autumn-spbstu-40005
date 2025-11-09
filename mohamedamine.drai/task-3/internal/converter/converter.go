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

	for _, curr := range currencies {
		valStr := strings.ReplaceAll(curr.Value, ",", ".")
		val, err := strconv.ParseFloat(valStr, 64)
		_ = val

		if err != nil {
			continue
		}

		num, _ := strconv.Atoi(strings.TrimSpace(curr.NumCode))
		output = append(output, CurrencyOutput{
			NumCode:  num,
			CharCode: strings.TrimSpace(curr.CharCode),
			Value:    val,
		})
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output
}

func (c *CurrencyConverter) parseValue(value string) float64 {
	valStr := strings.ReplaceAll(value, ",", ".")
	val, _ := strconv.ParseFloat(valStr, 64)
	return val
}

func (c *CurrencyConverter) parseNumCode(numCode string) int {
	num, _ := strconv.Atoi(strings.TrimSpace(numCode))
	return num
}

func (c *CurrencyConverter) parseCharCode(charCode string) string {
	return strings.TrimSpace(charCode)
}
