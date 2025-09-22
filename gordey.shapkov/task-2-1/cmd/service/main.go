package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MaxTemperature = 30
	MinTemperature = 15
)

var (
	errInput    = errors.New("invalid input")
	errOperator = errors.New("incorrect operator")
)

func changeTemperature(preferences []string) error {
	maxT := MaxTemperature
	minT := MinTemperature

	for _, pref := range preferences {
		parts := strings.Fields(pref)
		operator := parts[0]

		temp, err := strconv.Atoi(parts[1])
		if err != nil {
			return errInput
		}

		switch operator {
		case "<=":
			if temp < maxT {
				maxT = temp
			}
		case ">=":
			if temp > minT {
				minT = temp
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
	var countDepartments int
	if _, err := fmt.Scan(&countDepartments); err != nil || countDepartments < 1 {
		fmt.Println(errInput.Error())
		return
	}

	for d := 0; d < countDepartments; d++ {
		var countEmployees int
		if _, err := fmt.Scan(&countEmployees); err != nil || countEmployees < 1 {
			fmt.Println(errInput.Error())
			return
		}

		var preferences []string
		for e := 0; e < countEmployees; e++ {
			var operator string
			var temperature int
			if _, err := fmt.Scan(&operator, &temperature); err != nil {
				fmt.Println(errInput.Error())
				return
			}

			preferences = append(preferences, operator+" "+strconv.Itoa(temperature))
		}

		if err := changeTemperature(preferences); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
