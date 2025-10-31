package parserxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"alexandra.karnauhova/task-1/internal/data"
	"golang.org/x/net/html/charset"
)

func ParseXML(filename string) (*data.ValArray, error) {
	itData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(itData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var res data.ValArray

	err = decoder.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &res, nil
}

func getCharset(label string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(label)
	if encoding == nil {
		return nil, errors.New("unsupported charset")
	}

	return encoding.NewDecoder().Reader(input), nil
}
