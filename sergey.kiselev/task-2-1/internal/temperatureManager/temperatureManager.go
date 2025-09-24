package temperatureManager

import "fmt"

const (
	maxTemperature = 30
	minTemperature = 15
)

type TemperatureManager struct {
	maxTemp int
	minTemp int
}

func New() *TemperatureManager {
	return &TemperatureManager{
		maxTemp: maxTemperature,
		minTemp: minTemperature,
	}
}

func (temp *TemperatureManager) Update(operator string, temperature int) error {
	switch operator {
	case "<=":
		if temperature < temp.maxTemp {
			temp.maxTemp = temperature
		}
	case ">=":
		if temperature > temp.minTemp {
			temp.minTemp = temperature
		}
	default:
		return fmt.Errorf("incorrect operator")
	}

	return nil
}

func (temp *TemperatureManager) GetComfortTemp() int {
	if temp.minTemp <= temp.maxTemp {
		return temp.minTemp
	}

	return -1
}
