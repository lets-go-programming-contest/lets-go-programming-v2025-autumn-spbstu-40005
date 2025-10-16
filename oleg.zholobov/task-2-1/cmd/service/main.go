package main

import (
	"errors"
	"fmt"
)

func readInput() (int, error) {
	var value int
	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func readEmployeeRequest() (string, int, error) {
	var operator string
	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return "", 0, err
	}

	return operator, temp, nil
}

func processDepartment(employeeCount int) error {
	for range employeeCount {
		operator, temp, err := readEmployeeRequest()
		if err != nil {
			return errors.New("invalid employee request")
		}
		fmt.Println(operator)
		fmt.Println(temp)
	}
	return nil
}

func main() {
	departmentCount, err := readInput()
	if err != nil {
		fmt.Println(-1)
		return
	}

	for range departmentCount {
		employeeCount, err := readInput()
		if err != nil {
			fmt.Println(-1)
			return
		}
		if err = processDepartment(employeeCount); err != nil {
			fmt.Println(-1)
		}
	}
}
