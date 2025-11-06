package xmlparsing

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
	"gordey.shapkov/task-3/internal/models"
)

func ParseXMLFile(path string) (*models.ValCurs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs models.ValCurs
	if err = decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("cannot decode file: %w", err)
	}

	return &valCurs, nil
}
