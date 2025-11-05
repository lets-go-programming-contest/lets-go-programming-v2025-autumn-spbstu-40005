package xmlparser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
	"oleg.zholobov/task-3/internal/datamodels"
)

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, fmt.Errorf("unsupported charset: %s", charsetLabel)
	}

	return encoding.NewDecoder().Reader(input), nil
}

func ParseXML(filepath string) ([]datamodels.Valute, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: close XML file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = getCharset

	var valCurs datamodels.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	return valCurs.Valutes, nil
}
