package employee

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/temperature"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

func ProcessEmployee(countEmployees int) error {
	tempRange := &temperature.TemperatureRange{Min: minTemperature, Max: maxTemperature}

	for range countEmployees {
		icon, tempValue, err := temperature.ReadTemperature()
		if err != nil {
			return fmt.Errorf("read temperature: %w", err)
		}

		if !tempRange.IsValid() {
			fmt.Println(-1)

			continue
		}

		temperature.UpdateTemperature(icon, tempValue, tempRange)

		if tempRange.IsValid() {
			fmt.Println(tempRange.Min)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}
