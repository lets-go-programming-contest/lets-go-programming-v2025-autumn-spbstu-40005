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

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing XML file")
		}
	}()

	return parseXML(file)
}

func parseXML(reader io.Reader) (*ValCurs, error) {
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	valCurs := new(ValCurs)
	if err := decoder.Decode(valCurs); err != nil {
		return nil, fmt.Errorf("XML decoding error: %w", err)
	}

	return valCurs, nil
}
