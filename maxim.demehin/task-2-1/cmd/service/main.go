package main

import (
	"fmt"
)

const (
	lowerLimit = 15
	upperLimit = 30
)

func selectOptimalTemperature(workersCnt int) (int, error) {

}

func main() {
	var departsCount int

	_, err := fmt.Scan(&departsCount)
	if err != nil {
		fmt.Println("Invalid count of departments")
		return
	}

	var (
		workersCount int
		temperature  int
		cmpSign      string
	)

	for departsInd := 0; departsInd < departsCount; departsInd++ {
		_, err = fmt.Scan(&workersCount)
		if err != nil {
			fmt.Println("Invalid count of workers")
			return
		}
		for workersInd := 0; workersInd < workersCount; workersInd++ {

		}

	}

}
