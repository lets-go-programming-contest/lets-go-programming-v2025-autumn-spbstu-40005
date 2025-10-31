package parserxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"alexandra.karnauhova/task-3/internal/data"
	"golang.org/x/net/html/charset"
)

var ErrUnsupportCharset = errors.New("unsupported charset")

func ParseXML(filename string) (*data.ValArray, error) {
	itData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(itData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var res data.ValArray
	if err := decoder.Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &res, nil
}

func getCharset(label string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(label)
	if encoding == nil {
		return nil, ErrUnsupportCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}
