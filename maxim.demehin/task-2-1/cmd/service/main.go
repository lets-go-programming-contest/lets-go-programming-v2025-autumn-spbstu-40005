package main

import (
	"errors"
	"fmt"
)

var errCmpInput = errors.New("unknown comparator")

const (
	lowerLimit = 15
	upperLimit = 30
)

type TemperatureRange struct {
	lower int
	upper int
}

func newTemperatureRange(lower, upper int) *TemperatureRange {
	return &TemperatureRange{
		lower: lower,
		upper: upper,
	}
}

func (tr *TemperatureRange) changeRange(cmp string, curr int) error {
	switch cmp {
	case ">=":
		if curr > tr.lower {
			tr.lower = curr
		}
	case "<=":
		if curr < tr.upper {
			tr.upper = curr
		}
	default:
		return errCmpInput
	}

	return nil
}

func (tr *TemperatureRange) isValid() bool {
	return tr.lower <= tr.upper
}

func (tr *TemperatureRange) getResult() int {
	if !tr.isValid() {
		return -1
	}

	return tr.lower
}

func (tr *TemperatureRange) handleOptimalTemperature(cmp string, curr int) (int, error) {
	err := tr.changeRange(cmp, curr)
	if err != nil {
		return -1, err
	}

	return tr.getResult(), nil
}

func handleDepartmentTemperatures(workersCnt int) error {
	tempRange := newTemperatureRange(lowerLimit, upperLimit)

	var (
		temperature int
		cmpSign     string
	)

	for range workersCnt {
		_, err := fmt.Scan(&cmpSign, &temperature)
		if err != nil {
			return fmt.Errorf("failed reading of cmp and temperature: %w", err)
		}

		res, err := tempRange.handleOptimalTemperature(cmpSign, temperature)
		if err != nil {
			return fmt.Errorf("failed to process temperature: %w", err)
		}

		fmt.Println(res)
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
		fmt.Printf("failed to read departments count: %v\n", err)

		return
	}

	for range departsCount {
		_, err = fmt.Scan(&workersCount)
		if err != nil {
			fmt.Printf("invalid workers count: %v\n", err)

			return
		}

		err = handleDepartmentTemperatures(workersCount)
		if err != nil {
			fmt.Printf("failed processing department: %v\n", err)

			return
		}
	}
}
