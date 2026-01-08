package processor

import (
	"sort"

	"github.com/dmitei/task-3/internal/models"
)

func SortCurrenciesByExchangeRate(currencyContainer models.CurrencyContainer) ([]models.CurrencyInfo, error) {
	sortedCurrencyList := make(models.CurrencyList, len(currencyContainer.CurrencyList))
	copy(sortedCurrencyList, currencyContainer.CurrencyList)

	sort.Sort(sortedCurrencyList)

	return sortedCurrencyList, nil
}
