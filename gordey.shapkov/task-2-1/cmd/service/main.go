package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation           = errors.New("invalid operation")
	errInvalidNumberOfDepartments = errors.New("invalid number of departments")
	errOutOfRangeTemperature      = errors.New("preferred temperature out of allowed range")
	errMaxBelowMin                = errors.New("max temperature below current min")
)

const (
	MaxTemp = 30
	MinTemp = 15
)

type TemperaturePreference struct {
	maxTemp, minTemp int
}

func NewTemperaturePreference(maxTemp, minTemp int) *TemperaturePreference {
	return &TemperaturePreference{maxTemp, minTemp}
}

func (temp *TemperaturePreference) getOptimalTemp() (int, error) {
	if temp.minTemp > temp.maxTemp {
		return -1, fmt.Errorf("%w: current max %d < current min %d", errMaxBelowMin, temp.maxTemp, temp.minTemp)
	}

	return temp.minTemp, nil
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferredTemp int) error {
	if preferredTemp < temp.minTemp || preferredTemp > temp.maxTemp {
		return fmt.Errorf("%w: %d", errOutOfRangeTemperature, preferredTemp)
	}

	switch sign {
	case ">=":
		return temp.handleGreaterEqual(preferredTemp)
	case "<=":
		return temp.handleLessEqual(preferredTemp)
	default:
		return errInvalidOperation
	}
}

func (temp *TemperaturePreference) handleGreaterEqual(preferredTemp int) error {
	if preferredTemp > temp.minTemp {
		temp.minTemp = preferredTemp
	}

	return nil
}

func (temp *TemperaturePreference) handleLessEqual(preferredTemp int) error {
	if preferredTemp < temp.maxTemp {
		temp.maxTemp = preferredTemp
	}

	return nil
}

func main() {
	var numberOfDepartments int
	if _, err := fmt.Scan(&numberOfDepartments); err != nil {
		fmt.Println(errInvalidNumberOfDepartments)

		return
	}

	for range numberOfDepartments {
		var numberOfEmployees int

		_, err := fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("invalid number of employees: ", err)

			return
		}

		var (
			sign         string
			preferedTemp int
		)

		temp := NewTemperaturePreference(MaxTemp, MinTemp)

		for range numberOfEmployees {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				fmt.Println("invalid prefered temperature: ", err)

				return
			}

			err = temp.changeTemperature(sign, preferedTemp)
			if err != nil {
				if errors.Is(err, errInvalidOperation) {
					fmt.Println("invalid operation:", err)

					return
				}
				fmt.Println(-1)

				continue
			}

			optTemp, err := temp.getOptimalTemp()
			if err != nil {
				fmt.Println(optTemp)

				continue
			}

			fmt.Println(temp.minTemp)
		}
	}
}
