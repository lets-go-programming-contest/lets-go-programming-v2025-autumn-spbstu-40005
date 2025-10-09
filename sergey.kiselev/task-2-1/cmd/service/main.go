package main

import (
	"fmt"

	"sergey.kiselev/task-2-1/internal/employee"
)

func readEmployee() (*employee.Employee, error) {
	var (
		operator    string
		temperature int
	)

	if _, err := fmt.Scan(&operator, &temperature); err != nil {
		return nil, fmt.Errorf("failed to read temperature and operator: %w", err)
	}

	return employee.New(operator, temperature), nil
}

func processDepartment() error {
	var countEmployees uint

	if _, err := fmt.Scan(&countEmployees); err != nil {
		return fmt.Errorf("failed to read countEmployees: %w", err)
	}

	for range countEmployees {
		empl, err := readEmployee()
		if err != nil {
			return fmt.Errorf("error readEmployee : %w", err)
		}

		comfortTemp, err := empl.Process()
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
