package temperature

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInput = errors.New("invalid input")
	ErrIcon  = errors.New("invalid icon")
)

func ReadTemperature() (string, int, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return "", 0, err
	}

	runes := []rune(input)
	const minInputLength = 3
	if len(runes) < minInputLength {
		return "", 0, ErrInput
	}

	icon := string(runes[:2])

	if icon != ">=" && icon != "<=" {
		return "", 0, ErrIcon
	}

	num, err := strconv.Atoi(string(runes[2:]))
	if err != nil {
		return "", 0, ErrInput
	}

	const minTemp = 15
	const maxTemp = 30
	if num < minTemp || num > maxTemp {
		return "", 0, ErrInput
	}

	return icon, num, nil
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
