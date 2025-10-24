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

	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/models"
)

func Run(cfg *config.Config) error {

	xmlData, err := readXMLData(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error reading xml file: %w", err)
	}

	jsonValutes := transformAndSort(xmlData)

	err = writeJSONData(cfg.OutputFile, jsonValutes)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

func readXMLData(path string) (*models.XMLValCurs, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var valCurs models.XMLValCurs

	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}

func transformAndSort(xmlData *models.XMLValCurs) []models.JSONValute {

	jsonValutes := make([]models.JSONValute, 0, len(xmlData.Valutes))

	for _, xmlVal := range xmlData.Valutes {

		valStr := strings.Replace(xmlVal.Value, ",", ".", 1)

		valFloat, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}

		numCodeInt, _ := strconv.Atoi(xmlVal.NumCode)

		jsonValutes = append(jsonValutes, models.JSONValute{
			NumCode:  numCodeInt,
			CharCode: xmlVal.CharCode,
			Value:    valFloat,
		})
	}

	sort.Slice(jsonValutes, func(i, j int) bool {
		return jsonValutes[i].Value > jsonValutes[j].Value
	})

	return jsonValutes
}

func writeJSONData(path string, data []models.JSONValute) error {

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
