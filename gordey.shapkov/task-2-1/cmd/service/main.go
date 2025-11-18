package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation           = errors.New("invalid operation")
	errInvalidNumberOfDepartments = errors.New("invalid number of departments")
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

func (temp *TemperaturePreference) getOptimalTemp() int {
	if temp.minTemp > temp.maxTemp {
		return -1
	}

	return temp.minTemp
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferredTemp int) error {
	switch sign {
	case ">=":
		temp.handleGreaterEqual(preferredTemp)

		return nil
	case "<=":
		temp.handleLessEqual(preferredTemp)

		return nil
	default:
		return errInvalidOperation
	}
}

func (temp *TemperaturePreference) handleGreaterEqual(preferredTemp int) {
	temp.minTemp = max(temp.minTemp, preferredTemp)
}

func (temp *TemperaturePreference) handleLessEqual(preferredTemp int) {
	temp.maxTemp = min(temp.maxTemp, preferredTemp)
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
				fmt.Println("invalid preferred temperature or sign: ", err)

				return
			}

			err = temp.changeTemperature(sign, preferedTemp)
			if err != nil {
				if errors.Is(err, errInvalidOperation) {
					fmt.Println("invalid operation:", err)

					return
				}
			}

			fmt.Println(temp.getOptimalTemp())
		}
	}
}
