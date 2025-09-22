package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MaxTemp = 30
	MinTemp = 15
)

func changeTemperature(preferences []string) {
	maxTemp := MaxTemp
	minTemp := MinTemp
	currTemp := 0

	if len(preferences) > 0 {
		parts := strings.Fields(preferences[0])
		preferredTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if sign == ">=" {
			currTemp = preferredTemp
		} else {
			currTemp = MinTemp
		}

		if currTemp < MinTemp {
			currTemp = MinTemp
		}
		if currTemp > MaxTemp {
			currTemp = MaxTemp
		}

		if printOrReturn(currTemp, len(preferences)-1) {
			return
		}
	}

	for idx := 1; idx < len(preferences); idx++ {
		parts := strings.Fields(preferences[idx])
		preferredTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferredTemp < MinTemp || preferredTemp > MaxTemp {
			currTemp = -1
		} else {
			currTemp, minTemp, maxTemp = computeCurrTemp(sign, preferredTemp, currTemp, minTemp, maxTemp)
		}

		if printOrReturn(currTemp, len(preferences)-idx-1) {
			return
		}
	}
}

func computeCurrTemp(sign string, preferredTemp, currTemp, minTemp, maxTemp int) (int, int, int) {
	switch sign {
	case ">=":
		if preferredTemp > maxTemp {
			currTemp = -1
		} else if preferredTemp > currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp > minTemp {
			minTemp = preferredTemp
		}
	case "<=":
		if preferredTemp < minTemp {
			currTemp = -1
		} else if preferredTemp < currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp < maxTemp {
			maxTemp = preferredTemp
		}
	default:
		currTemp = -1
	}
	return currTemp, minTemp, maxTemp
}

func printOrReturn(currTemp int, remain int) bool {
	if currTemp == -1 {
		for idx := 0; idx < remain+1; idx++ {
			fmt.Println(currTemp)
		}
		return true
	}
	fmt.Println(currTemp)
	return false
}

func main() {
	var number int
	if _, err := fmt.Scan(&number); err != nil {
		return
	}

	for deptIdx := 0; deptIdx < number; deptIdx++ {
		var count int
		if _, err := fmt.Scan(&count); err != nil {
			return
		}

		preferences := make([]string, 0, count)
		for empIdx := 0; empIdx < count; empIdx++ {
			var sign string
			var preferredTemp int
			if _, err := fmt.Scan(&sign, &preferredTemp); err != nil {
				return
			}
			preferences = append(preferences, sign+" "+strconv.Itoa(preferredTemp))
		}
		changeTemperature(preferences)
	}
}

