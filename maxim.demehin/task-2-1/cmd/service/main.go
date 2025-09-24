package main

import (
	"errors"
	"fmt"
)

var (
	errCmpInput   = errors.New("Unknown comparator")
	errInput      = errors.New("Input error")
	errOutOfRange = errors.New("Temperature is out of range")
)

const (
	lowerLimit = 15
	upperLimit = 30
)

func handleOptimalTemperature(cmp string, curr int, lower, upper *int) error {
	// Если диапазон уже пуст, всегда возвращаем -1
	if *lower > *upper {
		fmt.Println(-1)
		return nil
	}

	switch cmp {
	case ">=":
		switch {
		case curr >= *lower && curr <= *upper:
			*lower = curr // Поднимаем нижнюю границу
			fmt.Println(*lower)
		case curr > *upper:
			// Невозможное условие - сужаем диапазон до пустого
			*lower = *upper + 1
			fmt.Println(-1)
		default: // curr < *lower
			// Условие уже выполняется, оставляем как есть
			fmt.Println(*lower)
		}
	case "<=":
		switch {
		case curr <= *upper && curr >= *lower:
			*upper = curr // Опускаем верхнюю границу
			fmt.Println(*lower)
		case curr < *lower:
			// Невозможное условие - сужаем диапазон до пустого
			*upper = *lower - 1
			fmt.Println(-1)
		default: // curr > *upper
			// Условие уже выполняется, оставляем как есть
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
