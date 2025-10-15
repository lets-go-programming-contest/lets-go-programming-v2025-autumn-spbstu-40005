package main

import (
	"fmt"
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

	for i := 0; i < nEmployees; i++ {
		err = execEmployee()
		if err != nil {
			fmt.Printf("Employee error: %v\n", err)
		}
	}
	return nil
}

func execEmployee() error {
	var (
		mode      string
		parameter int
	)
	_, err := fmt.Scan(&mode, &parameter)
	if err != nil {
		return fmt.Errorf("invalid employee command: %w", err)
	}

	return nil
}
