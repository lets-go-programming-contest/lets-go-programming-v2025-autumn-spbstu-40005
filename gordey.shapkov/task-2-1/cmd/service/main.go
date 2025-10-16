package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation           = errors.New("invalid operation")
	errInvalidNumberOfDepartments = errors.New("invalid number of departments")
)

type TemperaturePreference struct {
	maxTemp, minTemp int
}

func NewTemperaturePreference(maxTemp, minTemp int) *TemperaturePreference {
	return &TemperaturePreference{maxTemp, minTemp}
}

func (temp *TemperaturePreference) getMaxTemp() int {
	return temp.maxTemp
}

func (temp *TemperaturePreference) getMinTemp() int {
	return temp.minTemp
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferredTemp, currTemp int) (int, error) {
	if preferredTemp < MinTemp || preferredTemp > MaxTemp {
		return -1, fmt.Errorf("preferred temperature %d out of allowed range", preferredTemp)
	}

	switch sign {
	case ">=":
		return temp.handleGreaterEqual(preferredTemp, currTemp)
	case "<=":
		return temp.handleLessEqual(preferredTemp, currTemp)
	default:
		return -1, errInvalidOperation
	}
}

func (temp *TemperaturePreference) handleGreaterEqual(preferredTemp, currTemp int) (int, error) {
	if preferredTemp > temp.getMaxTemp() {
		return -1, fmt.Errorf("preferred %d > current max %d", preferredTemp, temp.maxTemp)
	}

	if preferredTemp > temp.getMinTemp() {
		temp.minTemp = preferredTemp
	}

	if preferredTemp > currTemp {
		currTemp = preferredTemp
	}

	return currTemp, nil
}

func (temp *TemperaturePreference) handleLessEqual(preferredTemp int, currTemp int) (int, error) {
	if preferredTemp < temp.getMinTemp() {
		return -1, fmt.Errorf("preferred %d < current min %d", preferredTemp, temp.minTemp)
	}

	if preferredTemp < temp.getMaxTemp() {
		temp.maxTemp = preferredTemp
	}

	if preferredTemp < currTemp {
		currTemp = preferredTemp
	}

	return currTemp, nil
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

		currTemp := MinTemp
		temp := NewTemperaturePreference(MaxTemp, MinTemp)

		for range numberOfEmployees {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}

			if currTemp == -1 {
				fmt.Println(currTemp)

				continue
			}

			currTemp, err = temp.changeTemperature(sign, preferedTemp, currTemp)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(currTemp)
		}
	}
}
