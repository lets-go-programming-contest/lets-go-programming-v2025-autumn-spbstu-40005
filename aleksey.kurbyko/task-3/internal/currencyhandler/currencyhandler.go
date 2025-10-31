package currencyhandler

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
)

type CurrencyList struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Items   []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func (list CurrencyList) Len() int {
	return len(list.Items)
}

func (list CurrencyList) Swap(i, j int) {
	list.Items[i], list.Items[j] = list.Items[j], list.Items[i]
}

func (list CurrencyList) Less(i, j int) bool {
	firstValue, err1 := strconv.ParseFloat(list.Items[i].Value, 64)
	secondValue, err2 := strconv.ParseFloat(list.Items[j].Value, 64)
	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
	return firstValue < secondValue
}

func SortCurrencies(list *CurrencyList) {
	for i := range list.Items {
		list.Items[i].Value = strings.ReplaceAll(strings.TrimSpace(list.Items[i].Value), ",", ".")
	}
	sort.Sort(sort.Reverse(list))
}
