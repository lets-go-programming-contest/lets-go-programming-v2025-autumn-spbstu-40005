package sorter

import (
	"sort"

	"alina.duhanina/task-3/internal/model"
)

type ByValueDesc []model.Valute

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

	return float64(a[i].Value) > float64(a[j].Value)
}

func ConvertAndSortCurrencies(valCurs *model.ValCurs) []model.Valute {
	valutes := make([]model.Valute, 0, len(valCurs.Valutes))
	valutes = append(valutes, valCurs.Valutes...)

	sort.Sort(ByValueDesc(valutes))

	return valutes
}
