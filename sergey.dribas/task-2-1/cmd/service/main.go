package main

import (
	"errors"
	"fmt"
)

var ErrInput = errors.New("invalid input")

const (
	constMaxtemp = 30
	constMintemp = 15
)

type Department struct {
	MinTemp int
	MaxTemp int
}

func NewDepartment() *Department {
	return &Department{constMintemp, constMaxtemp}
}

func (depart *Department) ProcessConstraint(operand string, temp int) error {
	switch operand {
	case ">=":
		depart.MinTemp = max(depart.MinTemp, temp)
	case "<=":
		depart.MaxTemp = min(depart.MaxTemp, temp)
	default:
		return ErrInput
	}

	return nil
}

func (depart *Department) GetCurrentTemp() int {
	if depart.MinTemp > depart.MaxTemp {
		return -1
	}

	return depart.MinTemp
}

func readConstraint() (string, int, error) {
	var (
		temp    int
		operand string
	)

	if _, err := fmt.Scan(&operand, &temp); err != nil {
		return "", 0, fmt.Errorf("error while reading operand + temperature: %w", err)
	}

	return operand, temp, nil
}

func processDepartment() error {
	var DepartmentSize int
	if _, err := fmt.Scan(&DepartmentSize); err != nil {
		return fmt.Errorf("error while reading department size: %w", err)
	}

	depart := NewDepartment()

	for range DepartmentSize {
		operand, temp, err := readConstraint()
		if err != nil {
			return err
		}

		if err := depart.ProcessConstraint(operand, temp); err != nil {
			return err
		}

		fmt.Println(depart.GetCurrentTemp())
	}

	return nil
}

func main() {
	var departmentCount int
	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Println("Error while department count input:", err)

		return
	}

	for range departmentCount {
		if err := processDepartment(); err != nil {
			fmt.Println("Error while process department:", err)

			return
		}
	}
}
