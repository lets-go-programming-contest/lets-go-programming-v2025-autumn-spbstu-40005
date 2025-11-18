package processor

import (
	"fmt"
	"os"
	"sort"

	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/models"
	"polina.gavrilova/task-3/internal/parser"
)

func Run(cfg *config.Config) error {
	xmlData, err := parser.ReadXMLData(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error reading xml file: %w", err)
	}

	valutes := transformAndSort(xmlData)

	dirPerm := os.FileMode(0o755)
	filePerm := os.FileMode(0o600)

	err = parser.WriteJSONData(cfg.OutputFile, valutes, dirPerm, filePerm)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

func transformAndSort(xmlData *models.ValCurs) []models.Valute {
	// Правильное копирование - создаем новый слайс и копируем элементы
	valutes := make([]models.Valute, len(xmlData.Valutes))
	for i, v := range xmlData.Valutes {
		valutes[i] = models.Valute{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.Value,
		}
	}

	sort.Slice(valutes, func(i, j int) bool {
		return float64(valutes[i].Value) > float64(valutes[j].Value)
	})

	return valutes
}
