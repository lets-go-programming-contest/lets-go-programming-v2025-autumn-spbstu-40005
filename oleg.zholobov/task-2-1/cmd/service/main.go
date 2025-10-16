package main

import (
	"errors"
	"fmt"
)

type TemperatureRange struct {
	Min int
	Max int
}

func NewTemperatureRange() TemperatureRange {
	return TemperatureRange{Min: 15, Max: 30}
}

func (tr *TemperatureRange) IsValid() bool {
	return tr.Min <= tr.Max
}

func (tr *TemperatureRange) GetOptimalTemperature() (int, error) {
	if !tr.IsValid() {
		return -1, errors.New("no valid temperature")
	}
	return tr.Min, nil
}

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
	tr := NewTemperatureRange()
	fmt.Println(tr.GetOptimalTemperature())
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
