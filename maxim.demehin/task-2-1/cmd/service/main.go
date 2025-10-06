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

func (tr *TemperatureRange) handleGreaterEqual(curr int) (int, bool) {
	if curr > tr.upper {
		tr.lower = tr.upper + 1

		return 0, false
	}

	if curr < tr.lower {
		return tr.lower, true
	}

	tr.lower = curr

	return tr.lower, true
}

func (tr *TemperatureRange) handleLessEqual(curr int) (int, bool) {
	if curr < tr.lower {
		tr.upper = tr.lower - 1

		return 0, false
	}

	if curr > tr.upper {
		return tr.lower, true
	}

	tr.upper = curr

	return tr.lower, true
}

func (tr *TemperatureRange) processCmp(cmp string, curr int) (int, bool, error) {
	switch cmp {
	case ">=":
		res, isValid := tr.handleGreaterEqual(curr)

		return res, isValid, nil
	case "<=":
		res, isValid := tr.handleLessEqual(curr)

		return res, isValid, nil
	default:
		return 0, false, errCmpInput
	}
}

func (tr *TemperatureRange) handleOptimalTemperature(cmp string, curr int) error {
	if !tr.isValid() {
		fmt.Println(-1)

		return nil
	}

	res, isValid, err := tr.processCmp(cmp, curr)
	if err != nil {
		return err
	}

	if isValid {
		fmt.Println(res)
	} else {
		fmt.Println(-1)
	}

	return nil
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
			return errInput
		}

		if temperature > upperLimit || temperature < lowerLimit {
			return errOutOfRange
		}

		err = tempRange.handleOptimalTemperature(cmpSign, temperature)
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
			fmt.Printf("invalid workers count: %v\n", err)

			return
		}

		err = handleDepartmentTemperatures(workersCount)
		if err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}
