package main

import (
	"errors"
	"fmt"
)

const (
	LessEqual    = false
	GreaterEqual = true
)

var (
	errInputFail          = errors.New("input error")
	errInvalidOperation   = errors.New("invalid operation")
	errInvalidTemperature = errors.New("invalid temperature")
)

func scanOperator() (bool, error) {
	var op string
	if _, err := fmt.Scan(&op); err != nil {
		return false, errInputFail
	}

	switch op {
	case "<=":
		return LessEqual, nil
	case ">=":
		return GreaterEqual, nil
	default:
		return false, errInvalidOperation
	}
}

func scanTemp() (int, error) {
	var t int
	if _, err := fmt.Scan(&t); err != nil || t > 30 || t < 15 {
		return 0, errInvalidTemperature
	}
	return t, nil
}

func processDepartment(k int) error {
	maxT := 30
	minT := 15
	for range k {
		op, err := scanOperator()
		if err != nil {
			return err
		}

		t, err := scanTemp()
		if err != nil {
			return err
		}

		if op == LessEqual {
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
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println(errInputFail.Error())
		return
	}

	for range n {
		var k int
		if _, err := fmt.Scan(&k); err != nil {
			fmt.Println(errInputFail.Error())
			return
		}

		if err := processDepartment(k); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
