package main

import (
	"errors"
	"fmt"
)

var (
	errCmpInput = errors.New("unknown comparator")
)

const (
	lowerLimit = 15
	upperLimit = 30
)

type TemperatureRange struct {
	lower int
	upper int
}

func newTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		lower: lowerLimit,
		upper: upperLimit,
	}
}

func (tr *TemperatureRange) isValid() bool {
	return tr.lower <= tr.upper
}

func (tr *TemperatureRange) handleGreaterEqual(curr int) int {
	if curr < tr.lower {
		return tr.lower
	}

	if curr > tr.upper {
		tr.lower = curr

		return -1
	}

	tr.lower = curr

	return tr.upper
}

func (tr *TemperatureRange) handleLessEqual(curr int) int {
	if curr > tr.upper {
		return tr.upper
	}

	if curr < tr.lower {
		tr.upper = curr

		return -1
	}

	tr.upper = curr

	return tr.lower
}

func (tr *TemperatureRange) processCmp(cmp string, curr int) (int, error) {
	switch cmp {
	case ">=":
		return tr.handleGreaterEqual(curr), nil
	case "<=":
		return tr.handleLessEqual(curr), nil
	default:
		return 0, errCmpInput
	}
}

func (tr *TemperatureRange) handleOptimalTemperature(cmp string, curr int) (int, error) {
	if !tr.isValid() {
		return -1, nil
	}

	res, err := tr.processCmp(cmp, curr)
	if err != nil {
		return -1, err
	}

	return res, nil
}

func handleDepartmentTemperatures(workersCnt int) error {
	tempRange := newTemperatureRange()

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
