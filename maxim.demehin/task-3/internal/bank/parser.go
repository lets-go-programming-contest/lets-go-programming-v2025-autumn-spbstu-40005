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

func ParseXML(path string) (*ValCurs, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = getCharset

	var (
		valCursRaw ValCursRaw
		valCurs    ValCurs
	)

	err = decoder.Decode(&valCursRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	err = valCurs.convertFromRaw(valCursRaw)
	if err != nil {
		return nil, err
	}

	valCurs.sortValutes()

	return &valCurs, nil
}

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, ErrUnsupportedCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}

func ParseToJSON(valutes []Valute) ([]ValuteToOut, error) {
	valutesToOut := make([]ValuteToOut, len(valutes))

	for index, valute := range valutes {
		valutesToOut[index] = ValuteToOut{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    valute.Value,
		}
	}

	return valutesToOut, nil
}

func SaveValutesToFile(valutes []ValuteToOut, output string) error {
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
