package sorter

import (
	"sort"

	"alexandra.karnauhova/task-3/internal/data"
)

func SortByValueDesc(valutes []data.Valute) []data.Valute {
	sorted := make([]data.Valute, len(valutes))
	copy(sorted, valutes)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}
