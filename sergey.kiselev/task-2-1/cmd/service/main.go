package main

import (
	"errors"
	"fmt"
)

const (
	maxTemperature = 30
	minTemperature = 15
)

var (
	errInput    = errors.New("invalid input")
	errOperator = errors.New("incorrect operator")
)

func processEmployees(countEmployees int) error {
	maxT := maxTemperature
	minT := minTemperature

	for range countEmployees {
		var (
			operator    string
			temperature int
		)

		if _, err := fmt.Scan(&operator, &temperature); err != nil {
			return errInput
		}

		switch operator {
		case "<=":
			if temperature < maxT {
				maxT = temperature
			}
		case ">=":
			if temperature > minT {
				minT = temperature
			}
		default:
			return errOperator
		}

		if minT <= maxT {
			fmt.Println(minT)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

func main() {
	var countDepartaments int

	if _, err := fmt.Scan(&countDepartaments); err != nil || countDepartaments < 1 {
		fmt.Println(errInput.Error())

		return
	}

	for range countDepartaments {
		var countEmployees int

		if _, err := fmt.Scan(&countEmployees); err != nil || countEmployees < 1 {
			fmt.Println(errInput.Error())

			return
		}

		if err := processEmployees(countEmployees); err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
