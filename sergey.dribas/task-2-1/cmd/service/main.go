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
		return "", 0, ErrInput
	}

	return operand, temp, nil
}

func processDepartment() error {
	var size int
	if _, err := fmt.Scan(&size); err != nil {
		return ErrInput
	}

	depart := NewDepartment()

	for range size {
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
	var size int
	if _, err := fmt.Scan(&size); err != nil {
		fmt.Println("Error:", err)

		return
	}

	for range size {
		if err := processDepartment(); err != nil {
			fmt.Println("Error:", err)

			return
		}
	}
}
