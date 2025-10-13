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

func (temp *TemperaturePreference) setMaxTemp(maxTemp int) {
	temp.maxTemp = maxTemp
}

func (temp *TemperaturePreference) getMaxTemp() int {
	return temp.maxTemp
}

func (temp *TemperaturePreference) setMinTemp(minTemp int) {
	temp.minTemp = minTemp
}

func (temp *TemperaturePreference) getMinTemp() int {
	return temp.minTemp
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferedTemp int) (int, error) {
	if preferedTemp < MinTemp || preferedTemp > MaxTemp {
		return -1, nil
	}

	var currTemp int

	switch sign {
	case ">=":
		temp.handleGreaterEqual(preferedTemp, &currTemp)
	case "<=":
		temp.handleLessEqual(preferedTemp, &currTemp)
	default:
		return 0, errInvalidOperation
	}

	return currTemp, nil
}

func (temp *TemperaturePreference) handleGreaterEqual(preferedTemp int, currTemp *int) {
	if preferedTemp > temp.maxTemp {
		*currTemp = -1
	}

	if preferedTemp > temp.getMinTemp() {
		temp.setMinTemp(preferedTemp)
	}

	if preferedTemp > *currTemp {
		*currTemp = preferedTemp
	}
}

func (temp *TemperaturePreference) handleLessEqual(preferedTemp int, currTemp *int) {
	if preferedTemp < temp.minTemp {
		*currTemp = -1
	}

	if preferedTemp < temp.getMaxTemp() {
		temp.setMaxTemp(preferedTemp)
	}

	if preferedTemp < *currTemp {
		*currTemp = preferedTemp
	}
}

const (
	MaxTemp = 30
	MinTemp = 15
)

func main() {
	var numberOfDepartments, numberOfEmployees int
	if _, err := fmt.Scan(&numberOfDepartments); err != nil {
		fmt.Println(errInvalidNumberOfDepartments)

		return
	}

	for range numberOfDepartments {
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

			currTemp, err = temp.changeTemperature(sign, preferedTemp)
			if err != nil {
				fmt.Println(err)

				break
			}

			fmt.Println(currTemp)
		}
	}
}
