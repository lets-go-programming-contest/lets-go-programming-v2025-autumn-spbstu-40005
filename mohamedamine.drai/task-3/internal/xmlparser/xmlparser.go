package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Value    Float64 `xml:"Value"`
}

type Float64 float64

func (f *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	if err := d.DecodeElement(&raw, &start); err != nil {
		return fmt.Errorf("decode xml element: %w", err)
	}

	raw = strings.ReplaceAll(raw, ",", ".")

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*f = Float64(val)

	return nil
}

func ReadXML(path string) (*ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var vc ValCurs
	if err := decoder.Decode(&vc); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &vc, nil
}
