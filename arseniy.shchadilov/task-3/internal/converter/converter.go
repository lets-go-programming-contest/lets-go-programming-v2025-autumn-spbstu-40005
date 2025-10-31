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
    jsonIndent      = "    "
    dirPermissions  = 0o755
    filePermissions = 0o644
)

func ParseXMLFile(filePath string, target interface{}) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open XML file: %w", err)
    }
    defer closeFile(file, filePath)

    decoder := xml.NewDecoder(file)
    decoder.CharsetReader = charset.NewReaderLabel

    if err := decoder.Decode(target); err != nil {
        return fmt.Errorf("failed to decode XML: %w", err)
    }

    return nil
}

func WriteToJSON(data interface{}, outputPath string) error {
    jsonData, err := json.MarshalIndent(data, jsonPrefix, jsonIndent)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }

    outputDir := filepath.Dir(outputPath)
    if err := os.MkdirAll(outputDir, dirPermissions); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    if err := os.WriteFile(outputPath, jsonData, filePermissions); err != nil {
        return fmt.Errorf("failed to write JSON file: %w", err)
    }

    return nil
}

func closeFile(file *os.File, filePath string) {
    if err := file.Close(); err != nil {
        fmt.Fprintf(os.Stderr, "Warning: failed to close file %s: %v\n", filePath, err)
    }
}
