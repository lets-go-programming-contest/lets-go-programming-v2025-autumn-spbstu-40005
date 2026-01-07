package main

import (
	"errors"
	"fmt"
)

const (
	TemperatureMin = 15
	TemperatureMax = 30

	CountMin = 1
	CountMax = 1000
)

var (
	ErrIncorrectSign        = errors.New("incorrect sign")
	ErrIncorrectBorder      = errors.New("incorrect border")
	ErrIncorrectDepartments = errors.New("incorrect amount of departments")
	ErrIncorrectEmployees   = errors.New("incorrect amount of employees")
)

type Bounds struct {
	minTemp  int
	maxTemp  int
	possible bool
}

func NewBounds() Bounds {
	return Bounds{
		minTemp:  TemperatureMin,
		maxTemp:  TemperatureMax,
		possible: true,
	}
}

func (b *Bounds) Apply(sign string, temperature int) error {
	if temperature < TemperatureMin || temperature > TemperatureMax {
		return ErrIncorrectBorder
	}

	switch sign {
	case ">=":
		if temperature > b.minTemp {
			b.minTemp = temperature
		}
	case "<=":
		if temperature < b.maxTemp {
			b.maxTemp = temperature
		}
	default:
		return ErrIncorrectSign
	}

	if b.minTemp > b.maxTemp {
		b.possible = false
	}

	return nil
}

func readInt() (int, error) {
	var value int

	if _, err := fmt.Scan(&value); err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	return value, nil
}

func readWorkerRequest() (string, int, error) {
	var (
		sign        string
		temperature int
	)

	if _, err := fmt.Scan(&sign, &temperature); err != nil {
		return "", 0, fmt.Errorf("scan worker request: %w", err)
	}

	return sign, temperature, nil
}

func processDepartment(employeeCount int) error {
	bounds := NewBounds()

	for range employeeCount {
		sign, temperature, err := readWorkerRequest()
		if err != nil {
			return err
		}

		if bounds.possible {
			if err := bounds.Apply(sign, temperature); err != nil {
				return err
			}
		}

		if bounds.possible {
			fmt.Println(bounds.minTemp)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

func run() error {
	departmentCount, err := readInt()
	if err != nil {
		return ErrIncorrectDepartments
	}

	if departmentCount < CountMin || departmentCount > CountMax {
		return ErrIncorrectDepartments
	}

	for range departmentCount {
		employeeCount, err := readInt()
		if err != nil {
			return ErrIncorrectEmployees
		}

		if employeeCount < CountMin || employeeCount > CountMax {
			return ErrIncorrectEmployees
		}

		if err := processDepartment(employeeCount); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		return
	}
}
