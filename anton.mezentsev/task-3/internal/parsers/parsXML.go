package parsers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

var (
	ErrOpeningFileXML = errors.New("error opening file")
	ErrClosingFileXML = errors.New("error closing file")
	ErrXMLDecoding    = errors.New("error xml decoding")
)

func ParseXML[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpeningFileXML, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(ErrClosingFileXML)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var result T
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrXMLDecoding, err)
	}

	return &result, nil
}
