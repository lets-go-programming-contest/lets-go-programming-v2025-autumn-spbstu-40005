package valute

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var errUnsupportedCharset = errors.New("unsupported charset")

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}
type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ValCursFromXML(path string) (*ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml file: %w", err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("%w: %s", errUnsupportedCharset, charset)
	}

	var exch ValCurs
	if err := decoder.Decode(&exch); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &exch, nil
}

func (v ValCurs) Len() int {
	return len(v.Valutes)
}

func (v ValCurs) Swap(i, j int) {
	if i < 0 || j < 0 || i >= v.Len() || j >= v.Len() {
		panic("index out of range")
	}

	v.Valutes[i], v.Valutes[j] = v.Valutes[j], v.Valutes[i]
}

func (v ValCurs) Less(left, right int) bool {
	if left < 0 || right < 0 || left >= v.Len() || right >= v.Len() {
		panic("index out of range")
	}

	leftOperand, _ := strconv.ParseFloat(strings.ReplaceAll(v.Valutes[left].Value, ",", "."), 32)
	rightOperand, _ := strconv.ParseFloat(strings.ReplaceAll(v.Valutes[right].Value, ",", "."), 32)

	return leftOperand > rightOperand
}
