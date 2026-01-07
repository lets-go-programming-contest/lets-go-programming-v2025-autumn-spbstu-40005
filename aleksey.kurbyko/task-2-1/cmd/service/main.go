package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp = 15
	MaxTemp = 30

	MinCount = 1
	MaxCount = 1000
)

var (
	ErrBadDeptCount   = errors.New("incorrect amount of departments")
	ErrBadWorkerCount = errors.New("incorrect amount of employees")
	ErrBadSign        = errors.New("incorrect sign")
	ErrBadTemp        = errors.New("incorrect border")
)

func updateBounds(minVal int, maxVal int, sign string, t int) (int, int, bool, error) {
	if t < MinTemp || t > MaxTemp {
		return minVal, maxVal, false, ErrBadTemp
	}

	switch sign {
	case ">=":
		if t > minVal {
			minVal = t
		}
	case "<=":
		if t < maxVal {
			maxVal = t
		}
	default:
		return minVal, maxVal, false, ErrBadSign
	}

	if minVal > maxVal {
		return minVal, maxVal, false, nil
	}

	return minVal, maxVal, true, nil
}

func main() {
	var deptCount int
	if _, err := fmt.Scan(&deptCount); err != nil || deptCount < MinCount || deptCount > MaxCount {
		fmt.Println("Error:", ErrBadDeptCount)
		return
	}

	for i := 0; i < deptCount; i++ {
		var workerCount int
		if _, err := fmt.Scan(&workerCount); err != nil || workerCount < MinCount || workerCount > MaxCount {
			fmt.Println("Error:", ErrBadWorkerCount)
			return
		}

		minVal := MinTemp
		maxVal := MaxTemp
		possible := true

		for j := 0; j < workerCount; j++ {
			var sign string
			var t int

			if _, err := fmt.Scan(&sign, &t); err != nil {
				fmt.Println("Error:", ErrBadTemp)
				return
			}

			if !possible {
				fmt.Println(-1)
				continue
			}

			var err error
			minVal, maxVal, possible, err = updateBounds(minVal, maxVal, sign, t)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if !possible {
				fmt.Println(-1)
			} else {
				fmt.Println(minVal)
			}
		}
	}
}
