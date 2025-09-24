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
	errArgument = errors.New("invalid argument")
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

func processEmployee(manager *TemperatureManager) error {
	var (
		operator    string
		temperature int
	)

	if _, err := fmt.Scan(&operator, &temperature); err != nil {
		return errInput
	}

	if err := manager.Update(operator, temperature); err != nil {
		return err
	}

	fmt.Println(manager.GetComfortTemp())

	return nil
}

func processDepartment() error {
	var countEmployees int

	if _, err := fmt.Scan(&countEmployees); err != nil {
		return errInput
	}

	if countEmployees < 1 {
		return errArgument
	}

	manager := TemperatureManager{maxTemperature, minTemperature}

	for range countEmployees {
		if err := processEmployee(&manager); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var countDepartaments int

	if _, err := fmt.Scan(&countDepartaments); err != nil {
		fmt.Println(errInput.Error())

		return
	}

	if countDepartaments < 1 {
		fmt.Print(errArgument.Error())

		return
	}

	for range countDepartaments {
		if err := processDepartment(); err != nil {
			fmt.Print(err.Error())

			return
		}
	}
}
