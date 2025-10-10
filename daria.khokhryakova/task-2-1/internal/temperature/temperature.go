package temperature

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrorInput = errors.New("Invalid input")
	ErrorIcon  = errors.New("Invalid icon")
)

func ReadTemperature() (string, int, error) {
	var input string
	fmt.Scan(&input)
	runes := []rune(input)
	if len(runes) < 3 {
		return "", 0, ErrorIcon
	}
	icon := string(runes[:2])
	num, _ := strconv.Atoi(string(runes[2:]))
	if icon != ">=" && icon != "<=" {
		return "", 0, ErrorIcon
	}
	if num < 15 || num > 30 {
		return "", 0, ErrorInput
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
		} else {
			return -1, false
		}
	case "<=":
		if temperature < *maxTemp {
			*maxTemp = temperature
		}
		if *minTemp <= *maxTemp {
			return *minTemp, true
		} else {
			return -1, false
		}
	default:
		return -1, false
	}
}
