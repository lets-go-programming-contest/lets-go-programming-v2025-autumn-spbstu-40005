package bank

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

func parseXML(path string) (*ValCurs, error) {
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
		return nil, fmt.Errorf("unsupported charset")
	}

	return encoding.NewDecoder().Reader(input), nil
}

func parseToJSON(valutes []Valute) ([]ValuteToOut, error) {

}
