package main

import (
	"errors"
	"fmt"
)

const (
	maxTemperature = 30
	minTemperature = 15
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
			return errors.New("invalid input for operator and temperature")
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
			return errors.New("incorrect operator")
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
		fmt.Println("Incorrect countDepartaments")

		return
	}

	for range countDepartaments {
		var countEmployees int

		if _, err := fmt.Scan(&countEmployees); err != nil || countEmployees < 1 {
			fmt.Println("Incorrect countEmployees")

			return
		}

		if err := processEmployees(countEmployees); err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
