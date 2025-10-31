package main

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/temperature"
)

const (
	minLimit = 15
	maxLimit = 30
)

func ReadTemperature() (string, int, error) {
	var operator string

	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return "", 0, fmt.Errorf("read temperature: %w", err)
	}

	return operator, temp, nil
}

func ProcessEmployee(countEmployees int) error {
	tempRange := temperature.NewTemperatureRange(minLimit, maxLimit)

	for range countEmployees {
		icon, tempValue, err := ReadTemperature()
		if err != nil {
			return fmt.Errorf("read temperature: %w", err)
		}

		err = temperature.UpdateTemperature(icon, tempValue, tempRange)
		if err != nil {
			return fmt.Errorf("update temperature: %w", err)
		}

		result := tempRange.GetResult()
		fmt.Println(result)
	}

	return nil
}

func ProcessDepartment(countDepartment int) error {
	for range countDepartment {
		err := ProcessSingleDepartment()
		if err != nil {
			return fmt.Errorf("process single department: %w", err)
		}
	}

	return nil
}

func ProcessSingleDepartment() error {
	var countEmployees int

	_, err := fmt.Scan(&countEmployees)
	if err != nil {
		return fmt.Errorf("read employee count: %w", err)
	}

	err = ProcessEmployee(countEmployees)
	if err != nil {
		return fmt.Errorf("process employee: %w", err)
	}

	return nil
}

func main() {
	var countDepartment int

	_, err := fmt.Scan(&countDepartment)
	if err != nil {
		fmt.Println("failed to read countDepartment:", err)

		return
	}

	err = ProcessDepartment(countDepartment)
	if err != nil {
		fmt.Println("failed in the function ProcessDepartment:", err)

		return
	}
}
