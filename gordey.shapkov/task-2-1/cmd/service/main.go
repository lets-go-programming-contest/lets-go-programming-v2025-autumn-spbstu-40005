package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MaxTemp = 30
	MinTemp = 15
)

var (
	errInput = errors.New("invalid input")
	errSign  = errors.New("incorrect sign")
)

func changeTemperature(preferences []string) error {
	currentMax := MaxTemp
	currentMin := MinTemp

	for idx := 0; idx < len(preferences); idx++ {
		parts := strings.Fields(preferences[idx])
		sign := parts[0]

		preferredTemp, err := strconv.Atoi(parts[1])
		if err != nil {
			return errInput
		}

		switch sign {
		case "<=":
			if preferredTemp < currentMax {
				currentMax = preferredTemp
			}
		case ">=":
			if preferredTemp > currentMin {
				currentMin = preferredTemp
			}
		default:
			return errSign
		}

		if currentMin <= currentMax {
			fmt.Println(currentMin)
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}

func main() {
	var numDepartments int

	if _, err := fmt.Scan(&numDepartments); err != nil || numDepartments < 1 {
		fmt.Println(errInput.Error())
		return
	}

	for deptIdx := 0; deptIdx < numDepartments; deptIdx++ {
		var numEmployees int

		if _, err := fmt.Scan(&numEmployees); err != nil || numEmployees < 1 {
			fmt.Println(errInput.Error())
			return
		}

		var preferences []string
		for empIdx := 0; empIdx < numEmployees; empIdx++ {
			var sign string
			var preferredTemp int

			if _, err := fmt.Scan(&sign, &preferredTemp); err != nil {
				fmt.Println(errInput.Error())
				return
			}

			preferences = append(preferences, sign+" "+strconv.Itoa(preferredTemp))
		}

		if err := changeTemperature(preferences); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

