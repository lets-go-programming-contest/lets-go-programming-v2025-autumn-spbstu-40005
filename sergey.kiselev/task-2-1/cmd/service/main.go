package main

import (
	"errors"
	"fmt"

	"sergey.kiselev/task-2-1/internal/temperatureManager"
)

var (
	errInput    = errors.New("invalid input")
	errArgument = errors.New("invalid argument")
)

func processEmployee(manager *temperatureManager.TemperatureManager) error {
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

	manager := temperatureManager.New()

	for range countEmployees {
		if err := processEmployee(manager); err != nil {
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
		fmt.Println(errArgument.Error())

		return
	}

	for range countDepartaments {
		if err := processDepartment(); err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
