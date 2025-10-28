package converter

import (
    "encoding/xml"
    "fmt"
    "os"
    "currency-converter/internal/currency"
    "golang.org/x/net/html/charset"
)

func ParseXMLFile(filePath string) ([]currency.Valute, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open XML file: %w", err)
    }
    defer file.Close()

    decoder := xml.NewDecoder(file)
    decoder.CharsetReader = charset.NewReaderLabel

    var valCurs currency.ValCurs
    if err := decoder.Decode(&valCurs); err != nil {
        return nil, fmt.Errorf("failed to decode XML: %w", err)
    }

    return valCurs.Valutes, nil
}
