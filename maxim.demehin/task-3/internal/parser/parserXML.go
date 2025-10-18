package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/TvoyBatyA12343/task-3/internal/bank"
	"golang.org/x/net/html/charset"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func ParseXML(path string) ([]bank.Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var valCurs bank.ValCurs

	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	valCurs.SortValutes()

	return valCurs.Valutes, nil
}

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, ErrUnsupportedCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}
