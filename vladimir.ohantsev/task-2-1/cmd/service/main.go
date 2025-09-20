package main

import (
	"errors"
	"fmt"
)

const (
	LessEqual = iota
	GreaterEqual
)

var (
	errInputFail          = errors.New("input error")
	errInvalidOperation   = errors.New("invalid operation")
	errInvalidTemperature = errors.New("invalid temperature")
)

func scanOperator() (int, error) {
	var op string
	if _, err := fmt.Scan(&op); err != nil {
		return 0, errInputFail
	}

	switch op {
	case "<=":
		return LessEqual, nil
	case ">=":
		return GreaterEqual, nil
	}

	return 0, errInvalidOperation
}

func scanTemp() (int, error) {
	var temp int
	if _, err := fmt.Scan(&temp); err != nil || temp > 30 || temp < 15 {
		return 0, errInvalidTemperature
	}

	return temp, nil
}

func processDepartment(countWorkers int) error {
	maxT := 30
	minT := 15

	for range countWorkers {
		operation, err := scanOperator()
		if err != nil {
			return err
		}

		t, err := scanTemp()
		if err != nil {
			return err
		}

		if operation == LessEqual {
			if maxT > t {
				maxT = t
			}
		} else {
			if minT < t {
				minT = t
			}
		}

		fmt.Println(curOptimum(minT, maxT))
	}

	return nil
}

func curOptimum(minT, maxT int) int {
	if maxT < minT {
		return -1
	}

	return minT
}

func main() {
	var countDepartments int
	if _, err := fmt.Scan(&countDepartments); err != nil {
		fmt.Println(errInputFail.Error())
		return
	}

	for range countDepartments {
		var countWorkers int
		if _, err := fmt.Scan(&countWorkers); err != nil {
			fmt.Println(errInputFail.Error())
			return
		}

		if err := processDepartment(countWorkers); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
