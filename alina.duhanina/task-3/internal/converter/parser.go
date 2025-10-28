package converter

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"alina.duhanina/task-3/internal/model"
)

func ParseXML(filePath string) (*model.ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open XML file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing XML file")
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs model.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("cannot decode XML: %w", err)
	}

	return &valCurs, nil
}

func parseValue(valueStr string) (float64, error) {
	normalized := strings.ReplaceAll(valueStr, ",", ".")

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("parse value %q: %w", valueStr, err)
	}

	return value, nil
}
