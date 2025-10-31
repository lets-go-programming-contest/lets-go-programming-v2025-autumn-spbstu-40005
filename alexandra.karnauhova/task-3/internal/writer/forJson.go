package writer

import (
	"encoding/json"
	"os"

	"alexandra.karnauhova/task-3/internal/data"
)

func CreateDirectory(directory string) error {
	return os.MkdirAll(directory, 0755)
}

func ParseJson(valutes []data.Valute, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(valutes)
}
