package currency

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func decodeWindows1251(data []byte) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	return decoder.Bytes(data)
}

func ParseXMLData(filePath string) (*ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	decodedData, err := decodeWindows1251(data)
	if err != nil {
		decodedData = data
	}

	var valCurs ValCurs
	err = xml.Unmarshal(decodedData, &valCurs)
	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}

func convertValue(valueStr string) (float64, error) {
	normalizStr := strings.Replace(valueStr, ",", ".", -1)
	value, err := strconv.ParseFloat(normalizStr, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func convertNumCode(numCodeStr string) (int, error) {
	numCode, err := strconv.Atoi(numCodeStr)
	if err != nil {
		return 0, err
	}

	return numCode, nil
}

func ProcessCurrencies(valCurs *ValCurs) ([]CurrencyResult, error) {
	var results []CurrencyResult
	for _, valute := range valCurs.Valutes {
		value, err := convertValue(valute.Value)
		if err != nil {
			return nil, err
		}
		numCode, err := convertNumCode(valute.NumCode)
		if err != nil {
			return nil, err
		}
		result := CurrencyResult{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func SaveResults(results []CurrencyResult, outputPath string) error {
	dir := filepath.Dir(outputPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(results)
	if err != nil {
		return err
	}
	return nil
}
