package parserxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

var ErrUnsupportCharset = errors.New("unsupported charset")

func ParseXML[T any](filename string) (*T, error) {
	itData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(itData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var res T
	if err := decoder.Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &res, nil
}

func getCharset(label string, input io.Reader) (io.Reader, error) {
	encoding, name := charset.Lookup(label)
	if encoding == nil || name == "" {
		return nil, ErrUnsupportCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}
