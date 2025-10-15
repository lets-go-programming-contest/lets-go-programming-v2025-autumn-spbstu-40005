package main

import (
	"fmt"

	"polina.gavrilova/task-2-1/internal/temperature"
)

func main() {
	var nDepartments int
	_, err := fmt.Scan(&nDepartments)
	if err != nil || nDepartments <= 0 {
		fmt.Printf("Invalid number of departments: %v\n", err)
		return
	}

	for i := 0; i < nDepartments; i++ {
		err = execDepartment()
		if err != nil {
			fmt.Printf("Department %d error: %v\n", i+1, err)
		}
	}
}

func execDepartment() error {
	var nEmployees int
	_, err := fmt.Scan(&nEmployees)
	if err != nil || nEmployees < 0 {
		return fmt.Errorf("invalid number of employees: %w", err)
	}

	tempCondition := &temperature.TempCondition{
		CurMin:  temperature.MinTemp,
		CurMax:  temperature.MaxTemp,
		CurTemp: temperature.MinTemp,
	}

	for i := 0; i < nEmployees; i++ {
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
		fmt.Println(-1)
		return nil
	}
	fmt.Println(tempCondition.CurTemp)
	return nil
}
