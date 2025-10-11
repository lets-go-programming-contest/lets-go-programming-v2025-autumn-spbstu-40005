package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXMLFile(filePath string) (*ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error when opening a file: %w", err)
	}
	defer file.Close()

	return ParseXML(file)
}

func ParseXML(reader io.Reader) (*ValCurs, error) {
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel
	valCurs := new(ValCurs)
	err := decoder.Decode(valCurs)
	if err != nil {
		return nil, fmt.Errorf("XML decoding error: %w", err)
	}

	return valCurs, nil
}
