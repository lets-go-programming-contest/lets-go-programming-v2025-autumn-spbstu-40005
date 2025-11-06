package main

import (
	"errors"
	"fmt"
)

type TemperatureProcessor struct {
	upper int
	lower int
}

const (
	UpperTemp = 30
	LowerTemp = 15
)

var (
	errOperator = errors.New("incorrect operator")
	errNotExist = errors.New("optimal temperature does not exist")
)

func NewTemperatureProcessor() *TemperatureProcessor {
	return &TemperatureProcessor{
		upper: UpperTemp,
		lower: LowerTemp,
	}
}

func (tp *TemperatureProcessor) apply(oper string, value int) error {
	switch oper {
	case "<=":
		if value < tp.upper {
			tp.upper = value
		}
	case ">=":
		if value > tp.lower {
			tp.lower = value
		}
	default:
		return errOperator
	}

	return nil
}

func (tp *TemperatureProcessor) getOptimalTemp() (int, error) {
	if tp.upper < tp.lower {
		return 0, errNotExist
	}

	return tp.lower, nil
}

func main() {
	var departmentCount int
	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Println("error reading department count: ", err)

		return
	}

	for range departmentCount {
		processor := NewTemperatureProcessor()

		var employeeCount int
		if _, err := fmt.Scan(&employeeCount); err != nil {
			fmt.Println("error reading employee count: ", err)

			return
		}

		for range employeeCount {
			var (
				oper  string
				value int
			)

			if _, err := fmt.Scan(&oper, &value); err != nil {
				fmt.Println("error reading new temperature: ", err)

				return
			}

			if err := processor.apply(oper, value); err != nil {
				fmt.Println("error applying operation: ", err)

				return
			}

			if res, err := processor.getOptimalTemp(); err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(res)
			}
		}
	}
}
