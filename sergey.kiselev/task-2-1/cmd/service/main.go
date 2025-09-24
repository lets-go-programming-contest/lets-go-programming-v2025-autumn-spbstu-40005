package main

import (
	"errors"
	"fmt"
)

const (
	maxTemperature = 30
	minTemperature = 15
)

var (
	errInput    = errors.New("invalid input")
	errOperator = errors.New("incorrect operator")
)

type TemperatureManager struct {
	maxTemp int
	minTemp int
}

func (temp *TemperatureManager) Update(operator string, temperature int) error {
	switch operator {
	case "<=":
		if temperature < temp.maxTemp {
			temp.maxTemp = temperature
		}
	case ">=":
		if temperature > temp.minTemp {
			temp.minTemp = temperature
		}
	default:
		return errOperator
	}
	return nil
}

func (temp *TemperatureManager) GetComfortTemp() int {
	if temp.minTemp <= temp.maxTemp {
		return temp.minTemp
	}
	return -1
}
func processEmployees(countEmployees int) error {
	manadger := TemperatureManager{maxTemperature, minTemperature}
	for range countEmployees {
		var (
			operator    string
			temperature int
		)

		if _, err := fmt.Scan(&operator, &temperature); err != nil {
			return errInput
		}

		if err := manadger.Update(operator, temperature); err != nil {
			return err
		}
		fmt.Println(manadger.GetComfortTemp())
	}

	return nil
}

func main() {
	var countDepartaments int

	if _, err := fmt.Scan(&countDepartaments); err != nil || countDepartaments < 1 {
		fmt.Println(errInput.Error())

		return
	}

	for range countDepartaments {
		var countEmployees int

		if _, err := fmt.Scan(&countEmployees); err != nil || countEmployees < 1 {
			fmt.Println(errInput.Error())

			return
		}

		if err := processEmployees(countEmployees); err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
