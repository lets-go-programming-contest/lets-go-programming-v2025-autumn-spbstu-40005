package temperature

import "errors"

var ErrIcon = errors.New("invalid icon")

type TemperatureRange struct {
	Min int
	Max int
}

func (temp *TemperatureRange) IsValid() bool {
	return temp.Min <= temp.Max
}

func UpdateTemperature(icon string, temperature int, tempRange *TemperatureRange) error {
	switch icon {
	case ">=":
		if temperature > tempRange.Min {
			tempRange.Min = temperature
		}

		return nil
	case "<=":
		if temperature < tempRange.Max {
			tempRange.Max = temperature
		}

		return nil
	default:
		return ErrIcon
	}
}
