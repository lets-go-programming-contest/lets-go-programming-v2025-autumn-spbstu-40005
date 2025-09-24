package main

import (
	"errors"
	"fmt"
)

var (
	errCmpInput   = errors.New("unknown comparator")
	errInput      = errors.New("input error")
	errOutOfRange = errors.New("temperature is out of range")
)

const (
	lowerLimit = 15
	upperLimit = 30
)

func handleOptimalTemperature(cmp string, curr int, lower, upper *int) error {
	if *lower > *upper {
		fmt.Println(-1)
		return nil
	}

	switch cmp {
	case ">=":
		switch {
		case curr >= *lower && curr <= *upper:
			*lower = curr
			fmt.Println(*lower)
		case curr > *upper:
			*lower = *upper + 1
			fmt.Println(-1)
		default:
			fmt.Println(*lower)
		}
	case "<=":
		switch {
		case curr <= *upper && curr >= *lower:
			*upper = curr
			fmt.Println(*lower)
		case curr < *lower:
			*upper = *lower - 1
			fmt.Println(-1)
		default:
			fmt.Println(*lower)
		}
	default:

		return errCmpInput
	}

	return nil
}

func handleDepartmentTemperatures(workersCnt int) error {
	var (
		currLower   = lowerLimit
		currUpper   = upperLimit
		temperature int
		cmpSign     string
	)

	for range workersCnt {
		_, err := fmt.Scan(&cmpSign, &temperature)
		if err != nil {

			return errInput
		}

		if temperature > upperLimit || temperature < lowerLimit {

			return errOutOfRange
		}

		err = handleOptimalTemperature(cmpSign, temperature, &currLower, &currUpper)
		if err != nil {

			return err
		}
	}

	return nil
}

func main() {
	var (
		departsCount int
		workersCount int
	)

	_, err := fmt.Scan(&departsCount)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	for range departsCount {
		_, err = fmt.Scan(&workersCount)
		if err != nil {
			fmt.Println(errInput.Error())

			return
		}

		err = handleDepartmentTemperatures(workersCount)
		if err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
