package temperature

import (
	"errors"
	"fmt"
)

var (
	ErrInput = errors.New("invalid input")
	ErrIcon  = errors.New("invalid icon")
)

func ReadTemperature() (string, int, error) {
	var operator string
	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return "", 0, fmt.Errorf("read temperature: %w", err)
	}

	if operator != ">=" && operator != "<=" {
		return "", 0, ErrIcon
	}

	const minTemp = 15
	const maxTemp = 30

	if temp < minTemp || temp > maxTemp {
		return "", 0, ErrInput
	}

	return operator, temp, nil
}

func PreferenceTemperature(icon string, temperature int, minTemp, maxTemp *int) (int, bool) {
	switch icon {
	case ">=":
		if temperature > *minTemp {
			*minTemp = temperature
		}

		if *minTemp <= *maxTemp {
			return *minTemp, true
		}

		return -1, false
	case "<=":
		if temperature < *maxTemp {
			*maxTemp = temperature
		}

		if *minTemp <= *maxTemp {
			return *minTemp, true
		}

		return -1, false
	default:
		return -1, false
	}
}
