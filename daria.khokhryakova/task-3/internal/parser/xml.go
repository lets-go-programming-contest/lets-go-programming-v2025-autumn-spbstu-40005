package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/DariaKhokhryakova/task-3/internal/models"
	"golang.org/x/net/html/charset"
)

func ParseXMLData(filePath string) (*models.ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs models.ValCurs

	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &valCurs, nil
}
