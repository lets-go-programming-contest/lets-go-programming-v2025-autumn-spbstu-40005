package converter

import (
	"alina.duhanina/task-3/internal/model"
	"sort"
)

type ByValueDesc []model.CurrencyResult

func (a ByValueDesc) Len() int { return len(a) }

func (a ByValueDesc) Swap(i, j int) {
	if i < 0 || i >= len(a) || j < 0 || j >= len(a) {
		panic("index out of range in Swap")
	}
	a[i], a[j] = a[j], a[i]
}

func (a ByValueDesc) Less(i, j int) bool {
	if i < 0 || i >= len(a) || j < 0 || j >= len(a) {
		panic("index out of range in Less")
	}
	return a[i].Value > a[j].Value
}

func ConvertAndSortCurrencies(valCurs *model.ValCurs) []model.CurrencyResult {
	var currencies []model.CurrencyResult

	for _, valute := range valCurs.Valutes {
		value, err := parseValue(valute.Value)
		if err != nil {
			continue
		}

		currency := model.CurrencyResult{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
		currencies = append(currencies, currency)
	}

	sort.Sort(ByValueDesc(currencies))

	return currencies
}
