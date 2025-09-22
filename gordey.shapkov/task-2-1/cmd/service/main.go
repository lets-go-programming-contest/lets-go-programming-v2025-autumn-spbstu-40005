package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var number, count int
	if _, err := fmt.Scan(&number); err != nil {
		return
	}

	for range number {
		_, err = fmt.Scan(&count)
		if err != nil {
			return
		}

		var (
			sign         string
			preferedTemp int
			preferences  []string
		)

		for range count {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}

			temp := sign + " " + strconv.Itoa(preferedTemp)
			preferences = append(preferences, temp)
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

	for idx := range preferences {
		parts := strings.Fields(preferences[idx])
		preferedTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferedTemp < MinTemp || preferedTemp > MaxTemp {
			currTemp = -1
		}

		if idx == 0 {
			currTemp = MinTemp
			if sign == ">=" {
				currTemp = preferedTemp
				minTemp = preferedTemp
			} else {
				maxTemp = preferedTemp
			}
		} else {
			switch sign {
			case ">=":
				currTemp = handleGreaterEqual(preferedTemp, currTemp, &minTemp, maxTemp)
			case "<=":
				currTemp = handleLessEqual(preferedTemp, currTemp, minTemp, &maxTemp)
			default:
				currTemp = -1
			}
		}

		if currTemp == -1 {
			printRemaining(currTemp, len(preferences)-idx)

			return
		}

		fmt.Println(currTemp)
	}
}

func handleGreaterEqual(preferedTemp, currTemp int, minTemp *int, maxTemp int) int {
	if preferedTemp > maxTemp {
		return -1
	}

	if preferedTemp > currTemp {
		currTemp = preferedTemp
	}

	if preferedTemp > *minTemp {
		*minTemp = preferedTemp
	}

	return currTemp
}

func handleLessEqual(preferedTemp, currTemp int, minTemp int, maxTemp *int) int {
	if preferedTemp < minTemp {
		return -1
	}

	if preferedTemp < currTemp {
		currTemp = preferedTemp
	}

	if preferedTemp < *maxTemp {
		*maxTemp = preferedTemp
	}

	return currTemp
}

func printRemaining(value, count int) {
	for range count {
		fmt.Println(value)
	}
}
