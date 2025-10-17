package temperature

import "errors"

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

	if signs == ">=" {
		if temp > temps.Min {
			temps.Min = temp
		}
		if temps.IdealTemperature < temp {
			temps.IdealTemperature = temp
		}
	} else if signs == "<=" {
		if temp < temps.Max {
			temps.Max = temp
		}
		if temps.IdealTemperature > temp {
			temps.IdealTemperature = temp
		}
	}

	if temps.Max < temps.Min {
		return errors.New("error in the input logic")
	}

	temps.Temps = append(temps.Temps, newTemp)

	return nil
}
