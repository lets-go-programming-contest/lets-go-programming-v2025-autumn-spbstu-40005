package main

import (
	"errors"
	"fmt"
)

var errIN = errors.New("invalid input")

func processDepartment() error {
	var departmentSize int

	_, err := fmt.Scan(&departmentSize)
	if err != nil {
		return errIN
	}

	minTemp := 15
	maxTemp := 30

	for range departmentSize {
		var (
			operand string
			temp    int
		)

		_, err := fmt.Scan(&operand, &temp)
		if err != nil {
			return errIN
		}

		switch operand {
		case ">=":
			if temp > minTemp {
				minTemp = temp
			}
		case "<=":
			if temp < maxTemp {
				maxTemp = temp
			}
		default:
			return errIN
		}

		if minTemp > maxTemp {
			fmt.Println("-1")
		} else {
			fmt.Println(minTemp)
		}
	}

	return nil
}

func main() {
	var size int
	if _, err := fmt.Scan(&size); err != nil {
		return
	}

	for range size {
		if processDepartment() != nil {
			return
		}
	}
}
