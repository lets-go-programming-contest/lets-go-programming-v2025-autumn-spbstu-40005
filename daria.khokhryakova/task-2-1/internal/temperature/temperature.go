package temperature

import "errors"

var ErrIcon = errors.New("invalid icon")

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange(minVal, maxVal int) *TemperatureRange {
	return &TemperatureRange{
		minTemp: minVal,
		maxTemp: maxVal,
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
