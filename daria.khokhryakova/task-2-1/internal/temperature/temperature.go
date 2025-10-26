package temperature

import "errors"

var (
	ErrIcon = errors.New("invalid icon")
	ErrTemp = errors.New("temperature out of range")
)

const (
	minLimit = 15
	maxLimit = 30
)

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		minTemp: minLimit,
		maxTemp: maxLimit,
	}
}

func (temp *TemperatureRange) IsValid() bool {
	return temp.minTemp <= temp.maxTemp
}

func (temp *TemperatureRange) GetResult() int {
	if !temp.IsValid() {
		return -1
	}

	return temp.minTemp
}

func UpdateTemperature(icon string, temperature int, tempRange *TemperatureRange) error {
	if temperature < minLimit || temperature > maxLimit {
		return ErrTemp
	}

	switch icon {
	case ">=":
		if temperature > tempRange.minTemp {
			tempRange.minTemp = temperature
		}

		return nil
	case "<=":
		if temperature < tempRange.maxTemp {
			tempRange.maxTemp = temperature
		}

		return nil
	default:
		return ErrIcon
	}
}
