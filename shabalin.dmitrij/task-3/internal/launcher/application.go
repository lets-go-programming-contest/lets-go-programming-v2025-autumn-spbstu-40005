package launcher

import (
	"fmt"
	"os"

	"github.com/dmitei/task-3/internal/configuration"
	"github.com/dmitei/task-3/internal/models"
	"github.com/dmitei/task-3/internal/processor"
	"github.com/dmitei/task-3/internal/storage"
	"github.com/dmitei/task-3/internal/xmlhandler"
)

func StartApplication(configurationFilePath string) error {
	appConfiguration, configLoadError := configuration.LoadApplicationConfig(configurationFilePath)
	if configLoadError != nil {
		return fmt.Errorf("failed to load configuration: %w", configLoadError)
	}

	sourceFileContent, fileReadError := os.ReadFile(appConfiguration.SourceFile)
	if fileReadError != nil {
		return fmt.Errorf("cannot read source file %q: %w", appConfiguration.SourceFile, fileReadError)
	}

	currencyContainer := models.CurrencyContainer{}

	if xmlParseError := xmlhandler.ParseXMLFile(sourceFileContent, &currencyContainer); xmlParseError != nil {
		return fmt.Errorf("failed to parse XML file: %w", xmlParseError)
	}

	sortedCurrencies, sortError := processor.SortCurrenciesByExchangeRate(currencyContainer)
	if sortError != nil {
		return fmt.Errorf("failed to sort currencies: %w", sortError)
	}

	if saveError := storage.SaveToJSONFile(appConfiguration.DestinationFile, sortedCurrencies); saveError != nil {
		return fmt.Errorf("failed to save results: %w", saveError)
	}

	return nil
}
