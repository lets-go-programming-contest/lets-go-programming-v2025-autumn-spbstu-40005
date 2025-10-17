package temperature

import (
	"errors"
)

var ErrInputLogic = errors.New("error in the input logic")

type RangeTemp struct {
	Range       string
	Temperature int
}

type TempManager struct {
	Temps            []RangeTemp
	IdealTemperature int
	Max, Min         int
}

func (temps *TempManager) AddTemp(signs string, temp int) error {
	newTemp := RangeTemp{
		Range:       signs,
		Temperature: temp,
	}

	switch signs {
	case ">=":
		temps.Min = max(temps.Min, temp)
	case "<=":
		temps.Max = min(temps.Max, temp)
	default:
		return ErrInputLogic
	}

	if temps.Min > temps.Max {
		temps.IdealTemperature = -1
	} else {
		temps.IdealTemperature = temps.Min
	}

	temps.Temps = append(temps.Temps, newTemp)

	return nil
}
