package datamodels

import (
	"sort"
	"strconv"
	"strings"
)

func ConvertAndSort(valutes []Valute) []Currency {
	currencies := make([]Currency, 0, len(valutes))

	for _, valute := range valutes {
		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			continue
		}

		valueStr := strings.Replace(valute.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		value = value / float64(valute.Nominal)

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies
}
