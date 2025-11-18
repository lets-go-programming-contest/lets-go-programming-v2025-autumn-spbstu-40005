package processor

import (
	"fmt"
	"os"
	"sort"

	"polina.gavrilova/task-3/internal/config"
	"polina.gavrilova/task-3/internal/models"
	"polina.gavrilova/task-3/internal/parser"
)

const (
	DefaultDirPermissions  = 0o755
	DefaultFilePermissions = 0o600
)

func Run(cfg *config.Config) error {
	xmlData, err := parser.ReadXMLData(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error reading xml file: %w", err)
	}

	valutes := transformAndSort(xmlData)

	dirPerm := getDirPermissions(cfg)
	filePerm := getFilePermissions(cfg)

	err = parser.WriteJSONData(cfg.OutputFile, valutes, dirPerm, filePerm)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

func getDirPermissions(cfg *config.Config) os.FileMode {
	if cfg.DirPerms != nil {
		perms := *cfg.DirPerms
		if perms < 0 {
			perms = 0
		}

		maskedPerms := perms & 0o777

		return os.FileMode(maskedPerms)
	}

	return DefaultDirPermissions
}

func getFilePermissions(cfg *config.Config) os.FileMode {
	if cfg.FilePerms != nil {
		perms := *cfg.FilePerms
		if perms < 0 {
			perms = 0
		}

		maskedPerms := perms & 0o777

		return os.FileMode(maskedPerms)
	}

	return DefaultFilePermissions
}

func transformAndSort(xmlData *models.ValCurs) []models.Valute {
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
