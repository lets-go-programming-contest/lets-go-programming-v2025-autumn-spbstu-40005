package bank

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func ParseXML(path string) ([]Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var valCurs ValCurs

	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	valCurs.sortValutes()

	return valCurs.Valutes, nil
}

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, ErrUnsupportedCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}

func SaveValutesToFile(valutes []Valute, output string) error {
	dir := filepath.Dir(output)

	err := os.MkdirAll(dir, 0)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(valutes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshal to JSON: %w", err)
	}

	err = os.WriteFile(output, jsonData, 0)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
