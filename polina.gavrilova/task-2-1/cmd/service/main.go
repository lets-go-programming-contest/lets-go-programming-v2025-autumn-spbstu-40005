package main

import (
	"fmt"

	"polina.gavrilova/task-2-1/internal/temperature"
)

func main() {
	var nDepartments int

	_, err := fmt.Scan(&nDepartments)
	if err != nil {
		fmt.Printf("Invalid number of departments: %v\n", err)

		return
	}

	for i := range nDepartments {
		err = execDepartment()
		if err != nil {
			fmt.Printf("Department %d error: %v\n", i+1, err)
		}
	}
}

func execDepartment() error {
	var nEmployees int

	_, err := fmt.Scan(&nEmployees)
	if err != nil {
		return fmt.Errorf("invalid number of employees: %w", err)
	}

	const (
		initialMinTemp = 15
		initialMaxTemp = 30
	)

	tempCondition := temperature.NewTempCondition(initialMinTemp, initialMaxTemp)

	for range nEmployees {
		err = execEmployee(tempCondition)
		if err != nil {
			fmt.Printf("Employee error: %v\n", err)
		}
	}

	return nil
}

func execEmployee(tempCondition *temperature.TempCondition) error {
	var (
		mode      string
		parameter int
	)

	_, err := fmt.Scan(&mode, &parameter)
	if err != nil {
		return fmt.Errorf("invalid employee command: %w", err)
	}

	err = tempCondition.Change(mode, parameter)
	if err != nil {
		return fmt.Errorf("invalid employee execution: %w", err)
	}

	if !tempCondition.HasValidRange() {
		fmt.Println(-1)

		return nil
	}

	minTemp, _ := tempCondition.GetCurrent()
	fmt.Println(minTemp)

	return nil
}
