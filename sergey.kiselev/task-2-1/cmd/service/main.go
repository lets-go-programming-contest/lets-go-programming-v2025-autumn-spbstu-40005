package main

import (
	"fmt"

	"sergey.kiselev/task-2-1/internal/temperaturemanager"
)

func processEmployee(manager *temperaturemanager.TemperatureManager) error {
	var (
		operator    string
		temperature int
	)

	if _, err := fmt.Scan(&operator, &temperature); err != nil {
		return fmt.Errorf("invalid input: %v", err)
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
		return fmt.Errorf("invalid input: %v", err)
	}

	if countEmployees < 1 {
		return fmt.Errorf("invalid argument: %v", countEmployees)
	}

	manager := temperaturemanager.New()

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
		fmt.Printf("failed to read department count: %v\n", countDepartaments)

		return
	}

	if countDepartaments < 1 {
		fmt.Println("department count must be positive")

		return
	}

	for range countDepartaments {
		if err := processDepartment(); err != nil {
			fmt.Printf("Error processing department: %v\n", err)

			return
		}
	}
}
