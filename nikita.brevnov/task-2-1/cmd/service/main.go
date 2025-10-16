package main

import "errors"

var ErrOperator = errors.New("invalid operator")

const (
	MinTemp = 15
	MaxTemp = 30
)

type climate struct {
	lowTemp  int
	HighTemp int
}

func climateController() *climate {
	return &climate{
		lowTemp:  MinTemp,
		HighTemp: MaxTemp,
	}
}

func main() {}
