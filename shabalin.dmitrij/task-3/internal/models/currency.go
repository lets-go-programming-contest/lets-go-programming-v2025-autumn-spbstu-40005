package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (currencyValue *CurrencyValue) UnmarshalXML(decoder *xml.Decoder, startElement xml.StartElement) error {
	var xmlStringValue string

	if decodeError := decoder.DecodeElement(&xmlStringValue, &startElement); decodeError != nil {
		return fmt.Errorf("cannot decode currency value element: %w", decodeError)
	}

	normalizedValue := strings.Replace(xmlStringValue, ",", ".", 1)

	parsedFloatValue, parseError := strconv.ParseFloat(normalizedValue, 64)
	if parseError != nil {
		return fmt.Errorf("cannot parse float value: %w", parseError)
	}

	*currencyValue = CurrencyValue(parsedFloatValue)

	return nil
}

type CurrencyInfo struct {
	NumericCode  int           `json:"num_code"  xml:"NumCode"`
	AlphaCode    string        `json:"char_code" xml:"CharCode"`
	ExchangeRate CurrencyValue `json:"value"     xml:"Value"`
}

type CurrencyContainer struct {
	ReportDate   string       `xml:"Date,attr"`
	ReportName   string       `xml:"name,attr"`
	CurrencyList CurrencyList `xml:"Valute"`
}

type CurrencyList []CurrencyInfo

func (currencyList CurrencyList) Len() int {
	return len(currencyList)
}

func (currencyList CurrencyList) Swap(firstIndex, secondIndex int) {
	currencyList[firstIndex], currencyList[secondIndex] = currencyList[secondIndex], currencyList[firstIndex]
}

func (currencyList CurrencyList) Less(firstIndex, secondIndex int) bool {
	return currencyList[firstIndex].ExchangeRate > currencyList[secondIndex].ExchangeRate
}
