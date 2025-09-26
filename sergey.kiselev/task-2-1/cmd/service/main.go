package main

import (
	"errors"
	"fmt"

	"sergey.kiselev/task-2-1/internal/temperature"
)

var errArgument = errors.New("invalid argument")

func processEmployee(manager *temperature.TemperatureManager) error {
	var (
		operator    string
		temperature int
	)

	if _, err := fmt.Scan(&operator, &temperature); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	if err := manager.Update(operator, temperature); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(manager.GetComfortTemp())

	return nil
}

func processDepartment() error {
	var countEmployees int

	if _, err := fmt.Scan(&countEmployees); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	if countEmployees < 1 {
		return fmt.Errorf("%w: employee count must be positive", errArgument)
	}

	manager := temperature.New()

	for range countEmployees {
		if err := processEmployee(manager); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func main() {
	var countDepartaments int

	if _, err := fmt.Scan(&countDepartaments); err != nil {
		fmt.Printf("failed to read department count: %s\n", err)

		return
	}

	if countDepartaments < 1 {
		fmt.Println("department count must be positive")

		return
	}

	for range countDepartaments {
		if err := processDepartment(); err != nil {
			fmt.Printf("Error processing department: %s\n", err)

			return
		}
	}
}
