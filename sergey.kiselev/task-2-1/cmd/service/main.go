package main

import (
	"fmt"

	"sergey.kiselev/task-2-1/internal/temperature"
)

func readEmployeeData() (string, int, error) {
	var (
		operator    string
		temperature int
	)

	if _, err := fmt.Scan(&operator, &temperature); err != nil {
		return "", 0, fmt.Errorf("failed to read temperature and operator: %w", err)
	}

	return operator, temperature, nil
}

func processDepartment() error {
	var countEmployees uint

	if _, err := fmt.Scan(&countEmployees); err != nil {
		return fmt.Errorf("failed to read countEmployees: %w", err)
	}

	manager := temperature.NewTemperatureManager()

	for range countEmployees {
		operator, temperature, err := readEmployeeData()
		if err != nil {
			return fmt.Errorf("error readEmployee : %w", err)
		}

		comfortTemp, err := manager.ProcessEmployee(operator, temperature)
		if err != nil {
			return fmt.Errorf("error process : %w", err)
		}

		fmt.Println(comfortTemp)
	}

	return nil
}

func main() {
	var countDepartaments uint

	if _, err := fmt.Scan(&countDepartaments); err != nil {
		fmt.Printf("failed to read department count: %s\n", err)

		return
	}

	for range countDepartaments {
		if err := processDepartment(); err != nil {
			fmt.Printf("Error processing department: %s\n", err)

			return
		}
	}
}
