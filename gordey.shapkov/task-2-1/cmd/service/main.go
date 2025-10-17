package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation           = errors.New("invalid operation")
	errInvalidNumberOfDepartments = errors.New("invalid number of departments")
	errOutOfRangeTemperature      = errors.New("preferred temperature out of allowed range")
	errPreferredAboveMax          = errors.New("preferred temperature above current max")
	errPreferredBelowMin          = errors.New("preferred temperature below current min")
)

type TemperaturePreference struct {
	maxTemp, minTemp int
}

func NewTemperaturePreference(maxTemp, minTemp int) *TemperaturePreference {
	return &TemperaturePreference{maxTemp, minTemp}
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferredTemp int) error {
	if preferredTemp < MinTemp || preferredTemp > MaxTemp {
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
	if preferredTemp > temp.maxTemp {
		return fmt.Errorf("%w: preferred %d > current max %d", errPreferredAboveMax, preferredTemp, temp.maxTemp)
	}

	if preferredTemp > temp.minTemp {
		temp.minTemp = preferredTemp
	}

	return nil
}

func (temp *TemperaturePreference) handleLessEqual(preferredTemp int) error {
	if preferredTemp < temp.minTemp {
		return fmt.Errorf("%w: preferred %d < current min %d", errPreferredBelowMin, preferredTemp, temp.minTemp)
	}

	if preferredTemp < temp.maxTemp {
		temp.maxTemp = preferredTemp
	}

	return nil
}

const (
	MaxTemp = 30
	MinTemp = 15
)

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
				return
			}

			err = temp.changeTemperature(sign, preferedTemp)
			if err != nil {
				switch {
				case errors.Is(err, errOutOfRangeTemperature),
					errors.Is(err, errPreferredAboveMax),
					errors.Is(err, errPreferredBelowMin):
					fmt.Println(-1)
				case errors.Is(err, errInvalidOperation):
					fmt.Println("invalid operation:", err)
				default:
					fmt.Println("unexpected error:", err)
				}

				continue
			}

			if temp.minTemp > temp.maxTemp {
				fmt.Println(-1)

				continue
			}

			fmt.Println(temp.minTemp)
		}
	}
}
