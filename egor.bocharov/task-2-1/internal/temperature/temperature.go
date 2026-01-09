package temperature

import (
	"errors"
	"fmt"
)

const (
	maxTemperature = 30
	minTemperature = 15
)

var errOperator = errors.New("incorrect operator")

type TemperatureManager struct {
	maxTemp int
	minTemp int
}

func NewTemperatureManager() *TemperatureManager {
	return &TemperatureManager{
		maxTemp: maxTemperature,
		minTemp: minTemperature,
	}
}

func (manager *TemperatureManager) Update(operator string, temperature int) error {
	switch operator {
	case "<=":
		if temperature < manager.maxTemp {
			manager.maxTemp = temperature
		}
	case ">=":
		if temperature > manager.minTemp {
			manager.minTemp = temperature
		}
	default:
		return errOperator
	}

	return nil
}

func (manager *TemperatureManager) GetComfortTemp() int {
	if manager.minTemp <= manager.maxTemp {
		return manager.minTemp
	}

	return -1
}

func (manager *TemperatureManager) ProcessEmployee(operator string, temperature int) (int, error) {
	if err := manager.Update(operator, temperature); err != nil {
		return 0, fmt.Errorf("error update temperature: %w", err)
	}

	return manager.GetComfortTemp(), nil
}
