package processor

import (
	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/models"
	"encoding/xml"
	"fmt"
	"os"
)

func Run(cfg *config.Config) error {

	xmlData, err := readXMLData(cfg.InputFile)
	if err != nil {

		return fmt.Errorf("error reading xml file: %w", err)
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
