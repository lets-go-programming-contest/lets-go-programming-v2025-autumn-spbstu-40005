package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var number, count int
	_, err := fmt.Scan(&number)
	if err != nil {
		return
	}

	for idx := 0; idx < number; idx++ {
		_, err = fmt.Scan(&count)
		if err != nil {
			return
		}

		var preferences []string
		for innerIdx := 0; innerIdx < count; innerIdx++ {
			var sign string
			var preferedTemp int
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}
			preferences = append(preferences, sign+" "+strconv.Itoa(preferedTemp))
		}

		changeTemperature(preferences)
	}
}

func changeTemperature(preferences []string) {
	const (
		MaxTemp = 30
		MinTemp = 15
	)

	maxTemp := MaxTemp
	minTemp := MinTemp
	currTemp := 0

	firstParts := strings.Fields(preferences[0])
	firstSign := firstParts[0]
	firstTemp, _ := strconv.Atoi(firstParts[1])

	if firstTemp < MinTemp || firstTemp > MaxTemp {
		currTemp = -1
		fmt.Println(currTemp)
		for j := 1; j < len(preferences); j++ {
			fmt.Println(currTemp)
		}
		fmt.Println()
		return
	}

	if firstSign == ">=" {
		currTemp = firstTemp
	} else {
		currTemp = MinTemp
	}
	fmt.Println(currTemp)

	for idx := 1; idx < len(preferences); idx++ {
		parts := strings.Fields(preferences[idx])
		sign := parts[0]
		preferedTemp, _ := strconv.Atoi(parts[1])

		currTemp = computeCurrTemp(currTemp, sign, preferedTemp, &minTemp, &maxTemp)
		if currTemp == -1 {
			for j := idx; j < len(preferences); j++ {
				fmt.Println(currTemp)
			}
			fmt.Println()
			return
		}
		fmt.Println(currTemp)
	}
}

func computeCurrTemp(currTemp int, sign string, preferedTemp int, minTemp *int, maxTemp *int) int {
	switch sign {
	case ">=":
		if preferedTemp > *maxTemp {
			return -1
		}
		if preferedTemp > currTemp {
			currTemp = preferedTemp
		}
		if preferedTemp > *minTemp {
			*minTemp = preferedTemp
		}
	case "<=":
		if preferedTemp < *minTemp {
			return -1
		}
		if preferedTemp < currTemp {
			currTemp = preferedTemp
		}
		if preferedTemp < *maxTemp {
			*maxTemp = preferedTemp
		}
	default:
		return -1
	}
	return currTemp
}

