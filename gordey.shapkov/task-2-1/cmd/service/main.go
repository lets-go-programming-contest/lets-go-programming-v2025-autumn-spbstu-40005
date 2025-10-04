package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation           = errors.New("invalid operation")
	errInvalidNumberOfEmployees   = errors.New("invalid number of employees")
	errInvalidNumberOfDepartments = errors.New("invalid number of departments")
)

type TemperaturePreference struct {
	maxTemp, minTemp, currTemp int
}

func NewTemperaturePreference(maxTemp, minTemp, currTemp int) *TemperaturePreference {
	return &TemperaturePreference{maxTemp, minTemp, currTemp}
}

func (temp *TemperaturePreference) setMaxTemp(maxTemp int) {
	temp.maxTemp = maxTemp
}

func (temp *TemperaturePreference) setMinTemp(minTemp int) {
	temp.minTemp = minTemp
}

func (temp *TemperaturePreference) setCurrTemp(currTemp int) {
	temp.currTemp = currTemp
}

func (temp *TemperaturePreference) getCurrTemp() int {
	return temp.currTemp
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
			fmt.Println(errInvalidNumberOfEmployees)

			return
		}

		var (
			sign         string
			preferedTemp int
		)

		temp := NewTemperaturePreference(MaxTemp, MinTemp, MinTemp)

		for range numberOfEmployees {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}

			if temp.currTemp == -1 {
				fmt.Println(temp.currTemp)

				continue
			}

			err := temp.changeTemperature(sign, preferedTemp)

			if err != nil {
				fmt.Println(err)

				break
			}

			fmt.Println(temp.getCurrTemp())
		}
	}
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferedTemp int) error {
	if preferedTemp < MinTemp || preferedTemp > MaxTemp {
		temp.setCurrTemp(-1)
	}

	switch sign {
	case ">=":
		handleGreaterEqual(preferedTemp, temp)
	case "<=":
		handleLessEqual(preferedTemp, temp)
	default:
		return errInvalidOperation
	}

	return nil
}

func handleGreaterEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp > temp.maxTemp {
		temp.setCurrTemp(-1)

		return
	}

	if preferedTemp > temp.currTemp {
		temp.setCurrTemp(preferedTemp)
	}

	if preferedTemp > temp.minTemp {
		temp.setMinTemp(preferedTemp)
	}
}

func handleLessEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp < temp.minTemp {
		temp.setCurrTemp(-1)

		return
	}

	if preferedTemp < temp.currTemp {
		temp.setCurrTemp(preferedTemp)
	}

	if preferedTemp < temp.maxTemp {
		temp.setMaxTemp(preferedTemp)
	}
}
