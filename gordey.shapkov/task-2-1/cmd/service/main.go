package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MaxTemp = 30
	MinTemp = 15
)

var (
	errInput = errors.New("invalid input")
	errSign  = errors.New("invalid operator")
)

func computeCurrTemp(currTemp, minTemp, maxTemp, preferredTemp int, sign string) int {
	switch sign {
	case ">=":
		if preferredTemp > maxTemp {
			return -1
		} else if preferredTemp > currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp > minTemp {
			minTemp = preferredTemp
		}
	case "<=":
		if preferredTemp < minTemp {
			return -1
		} else if preferredTemp < currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp < maxTemp {
			maxTemp = preferredTemp
		}
	default:
		return -1
	}
	return currTemp
}

func changeTemperature(preferences []string) {
	maxTemp := MaxTemp
	minTemp := MinTemp
	currTemp := 0

	if len(preferences) > 0 {
		parts := strings.Fields(preferences[0])
		preferredTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferredTemp < MinTemp || preferredTemp > MaxTemp {
			currTemp = -1
		} else {
			if sign == ">=" {
				currTemp = preferredTemp
			} else {
				currTemp = MinTemp
			}
		}
		fmt.Println(currTemp)
		minTemp = preferredTemp
		maxTemp = preferredTemp
	}

	for idx := 1; idx < len(preferences); idx++ {
		parts := strings.Fields(preferences[idx])
		preferredTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferredTemp < MinTemp || preferredTemp > MaxTemp {
			currTemp = -1
		} else {
			currTemp = computeCurrTemp(currTemp, minTemp, maxTemp, preferredTemp, sign)
		}

		if currTemp == -1 {
			for empIdx := idx; empIdx < len(preferences); empIdx++ {
				fmt.Println(currTemp)
			}
			return
		}
		fmt.Println(currTemp)
	}
}

func main() {
	var countDeparts int

	if _, err := fmt.Scan(&countDeparts); err != nil || countDeparts < 1 {
		fmt.Println(errInput.Error())
		return
	}

	for departIdx := 0; departIdx < countDeparts; departIdx++ {
		var countWorkers int
		if _, err := fmt.Scan(&countWorkers); err != nil || countWorkers < 1 {
			fmt.Println(errInput.Error())
			return
		}

		prefs := make([]string, 0, countWorkers)

		for workerIdx := 0; workerIdx < countWorkers; workerIdx++ {
			var sign string
			var preferredTemp int
			if _, err := fmt.Scan(&sign, &preferredTemp); err != nil {
				fmt.Println(errInput.Error())
				return
			}

			temp := sign + " " + strconv.Itoa(preferredTemp)
			prefs = append(prefs, temp)
		}
		changeTemperature(prefs)
	}
}
