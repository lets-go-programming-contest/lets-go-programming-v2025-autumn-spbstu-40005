package processor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/models"
)

func Run(cfg *config.Config) error {
	xmlData, err := readXMLData(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error reading xml file: %w", err)
	}

	jsonValutes, err := transformAndSort(xmlData)
	if err != nil {
		return err
	}

	err = writeJSONData(cfg.OutputFile, jsonValutes)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

func readXMLData(path string) (*models.XMLValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs models.XMLValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &valCurs, nil
}

func transformAndSort(xmlData *models.XMLValCurs) ([]models.JSONValute, error) {
	jsonValutes := make([]models.JSONValute, 0, len(xmlData.Valutes))

	for _, xmlVal := range xmlData.Valutes {
		if xmlVal.NumCode == "" || xmlVal.CharCode == "" || xmlVal.Value == "" {
			continue
		}

		valStr := strings.Replace(xmlVal.Value, ",", ".", 1)
		valFloat, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value for currency %s: %w", xmlVal.CharCode, err)
		}

		numCodeInt, err := strconv.Atoi(xmlVal.NumCode)
		if err != nil {
			return nil, fmt.Errorf("error parsing numcode for currency %s: %w", xmlVal.CharCode, err)
		}

		jsonValutes = append(jsonValutes, models.JSONValute{
			NumCode:  numCodeInt,
			CharCode: xmlVal.CharCode,
			Value:    valFloat,
		})
	}

	sort.Slice(jsonValutes, func(i, j int) bool {
		return jsonValutes[i].Value > jsonValutes[j].Value
	})

	return jsonValutes, nil
}

func writeJSONData(path string, data []models.JSONValute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
