package main

import (
	"errors"
	"fmt"
)

var (
	cmpInputError   = errors.New("Unknown comparator")
	inputError      = errors.New("Input error")
	outOfRangeError = errors.New("Temperature is out of range")
)

const (
	lowerLimit = 15
	upperLimit = 30
)

func printOptimalTemperature(cmp string, curr int, lower, upper *int) error {
	switch cmp {
	case ">=":
		if curr > *lower && curr <= *upper {
			*lower = curr
			fmt.Println(*lower)
		} else {
			fmt.Println(-1)
		}
	case "<=":
		if curr < *upper && curr >= *lower {
			fmt.Println(*lower)
		} else {
			fmt.Println(-1)
		}
	default:
		return cmpInputError
	}
	return nil
}

func handleTemperatures(workersCnt int) error {
	var (
		currLower   = lowerLimit
		currUpper   = upperLimit
		temperature int
		cmpSign     string
	)

	for workersInd := 0; workersInd < workersCnt; workersInd++ {
		_, err := fmt.Scan(&cmpSign)
		if err != nil {
			return inputError
		}
		_, err = fmt.Scan(&temperature)
		if err != nil {
			return inputError
		}
		if temperature > upperLimit || temperature < lowerLimit {
			return outOfRangeError
		}

		err = printOptimalTemperature(cmpSign, temperature, &currLower, &currUpper)
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
		fmt.Println(inputError.Error())
		return
	}

	for range departsCount {
		_, err = fmt.Scan(&workersCount)
		if err != nil {
			fmt.Println(inputError.Error())
			return
		}

		err = handleTemperatures(workersCount)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
}
