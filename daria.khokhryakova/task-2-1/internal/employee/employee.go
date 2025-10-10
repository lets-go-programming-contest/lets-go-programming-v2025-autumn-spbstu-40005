package employee

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/temperature"
)

func ProcessEmployee(countEmployees int) ([]int, error) {
	minTemp := 15
	maxTemp := 30
	results := []int{}
	processValid := true

	for range countEmployees {
		icon, tempValue, err := temperature.ReadTemperature()
		if err != nil {
			return nil, fmt.Errorf("read temperature: %w", err)
		}

		if !processValid {
			results = append(results, -1)
			continue
		}

		result, valid := temperature.PreferenceTemperature(icon, tempValue, &minTemp, &maxTemp)
		results = append(results, result)

		if !valid {
			processValid = false
		}
	}

	return results, nil
}
