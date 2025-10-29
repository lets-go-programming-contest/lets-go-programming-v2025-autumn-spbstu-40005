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

var ErrUnsupportedCharset = errors.New("unsupported charset")

func ParseXML[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var res T

	err = decoder.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &res, nil
}

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, ErrUnsupportedCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}
