package temperature

import (
	"errors"
	"fmt"
)

var (
	ErrInput = errors.New("invalid input")
	ErrIcon  = errors.New("invalid icon")
)

const (
	minTemperature = 15
	maxTemperature = 30
)

type TemperatureRange struct {
	Min int
	Max int
}

func (temp *TemperatureRange) IsValid() bool {
	return temp.Min <= temp.Max
}

func ReadTemperature() (string, int, error) {
	var operator string

	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return "", 0, fmt.Errorf("read temperature: %w", err)
	}

	switch operator {
	case ">=", "<=":
	default:
		return "", 0, ErrIcon
	}

	if temp < minTemperature || temp > maxTemperature {
		return "", 0, ErrInput
	}

	return operator, temp, nil
}

func UpdateTemperature(icon string, temperature int, tempRange *TemperatureRange) {
	switch icon {
	case ">=":
		if temperature > tempRange.Min {
			tempRange.Min = temperature
		}
	case "<=":
		if temperature < tempRange.Max {
			tempRange.Max = temperature
		}
	}
}
