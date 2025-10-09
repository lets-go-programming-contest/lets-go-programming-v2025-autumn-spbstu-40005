package employee

import (
	"errors"
	"fmt"
)

const (
	maxTemperature = 30
	minTemperature = 15
)

var errOperator = errors.New("incorrect operator")

type Employee struct {
	operator    string
	temperature int
	maxTemp     int
	minTemp     int
}

func New(oper string, temp int) *Employee {
	return &Employee{
		operator:    oper,
		temperature: temp,
		maxTemp:     maxTemperature,
		minTemp:     minTemperature,
	}
}

func (e *Employee) UpdateTemperature() error {
	switch e.operator {
	case "<=":
		if e.temperature < e.maxTemp {
			e.maxTemp = e.temperature
		}
	case ">=":
		if e.temperature > e.minTemp {
			e.minTemp = e.temperature
		}
	default:
		return errOperator
	}

	return nil
}

func (e *Employee) GetComfortTemp() int {
	if e.minTemp <= e.maxTemp {
		return e.minTemp
	}
	return -1
}

func (e *Employee) Process() (int, error) {
	if err := e.UpdateTemperature(); err != nil {
		return 0, fmt.Errorf("error update temperature: %w", err)
	}

	return e.GetComfortTemp(), nil
}
