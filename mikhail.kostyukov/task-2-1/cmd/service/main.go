package main

import (
	"errors"
	"fmt"
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
)

var ErrInvalidOperator = errors.New("invalid operator")

type TemperatureManager struct {
	minTemp int
	maxTemp int
}

func NewTemperatureManager() *TemperatureManager{
	return &TemperatureManager{
		minTemp: InitialMinTemp,
		maxTemp: InitialMaxTemp,
	}
}

func (tm *TemperatureManager) Update(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tm.minTemp {
			tm.minTemp = temp
		}
	case "<=":
		if temp < tm.maxTemp {
			tm.maxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func main() {
}

