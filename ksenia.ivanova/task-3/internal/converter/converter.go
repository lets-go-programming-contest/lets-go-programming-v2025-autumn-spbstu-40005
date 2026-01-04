// internal/converter/parse.go
package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

const (
	jsonPrefix      = ""
	jsonIndent      = "\t"
	dirPermissions  = 0o755
	filePermissions = 0o644
)

func ParseXMLFile(inputFile string, target interface{}) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed unmarshaling XML: %w", err)
	}

	return nil
}

func WriteToJSON(data interface{}, outputFile string) error {
	jsonData, err := json.MarshalIndent(data, jsonPrefix, jsonIndent)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %w", err)
	}

	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, dirPermissions); err != nil {
		return fmt.Errorf("failed creating directory %s: %w", outputDir, err)
	}

	if err := os.WriteFile(outputFile, jsonData, filePermissions); err != nil {
		return fmt.Errorf("failed writing file %s: %w", outputFile, err)
	}

	return nil
}
