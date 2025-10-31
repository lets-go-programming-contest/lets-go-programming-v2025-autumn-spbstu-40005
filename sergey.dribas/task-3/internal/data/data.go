package valute

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

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
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
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

func (v ValCurs) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= v.Len() || j >= v.Len() {
		panic("index out of range")
	}

	leftOperand, _ := strconv.ParseFloat(strings.ReplaceAll(v.Valutes[i].Value, ",", "."), 32)
	rightOperand, _ := strconv.ParseFloat(strings.ReplaceAll(v.Valutes[j].Value, ",", "."), 32)
	return leftOperand > rightOperand
}
