package bank

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Rate     string `xml:"VunitRate"`
}

func Parse(r io.Reader) (*Bank, error) {
	decoder := xml.NewDecoder(r)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	bank := new(Bank)
	if err := decoder.Decode(&bank); err != nil {
		return nil, fmt.Errorf("decoding currency bank: %w", err)
	}

	return bank, nil
}

func ParseFile(path string) (*Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}
	defer file.Close()

	return Parse(file)
}

func (b outputBank) sortByValueDown() {
	sort.Slice(
		b,
		func(i, j int) bool {
			return b[i].Value > b[j].Value
		},
	)
}

type outputCurrency struct {
	NumCode  int     `json:"num-code"`
	CharCode string  `json:"char-code"`
	Value    float64 `json:"value"`
}

type outputBank []outputCurrency

func (b *Bank) EncodeJSON(writer io.Writer) error {
	out := make(outputBank, len(b.Currencies))

	for index, currency := range b.Currencies {
		val, err := strconv.ParseFloat(strings.Replace(currency.Value, ",", ".", 1), 64)
		if err != nil {
			return fmt.Errorf("invalid type of value: %w", err)
		}

		out[index] = outputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    val,
		}
	}

	out.sortByValueDown()
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(&out); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
	}

	return nil
}

func (b *Bank) EncodeJSONToFIle(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer file.Close()

	return b.EncodeJSON(file)
}
