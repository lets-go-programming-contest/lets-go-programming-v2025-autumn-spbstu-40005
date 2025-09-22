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

	for range number {
		_, err = fmt.Scan(&count)
		if err != nil {
			return
		}

		var preferences []string
		for range count {
			var sign string
			var preferredTemp int
			_, err = fmt.Scan(&sign, &preferredTemp)
			if err != nil {
				return
			}
			temp := sign + " " + strconv.Itoa(preferredTemp)
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

	var (
		currTemp int
		maxTemp  = MaxTemp
		minTemp  = MinTemp
	)

	for idx := range preferences {
		parts := strings.Fields(preferences[idx])
		preferredTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferredTemp < MinTemp || preferredTemp > MaxTemp {
			currTemp = -1
		} else {
			if idx == 0 {
				if sign == ">=" {
					currTemp = preferredTemp
				} else {
					currTemp = MinTemp
				}
			}

			currTemp = computeCurrTemp(preferredTemp, sign, currTemp, &minTemp, &maxTemp)
		}

		if currTemp == -1 {

			for j := idx; j < len(preferences); j++ {
				fmt.Println(currTemp)
			}

			return
		}

		fmt.Println(currTemp)
	}
}

func computeCurrTemp(preferredTemp int, sign string, currTemp int, minTemp *int, maxTemp *int) int {
	switch sign {
	case ">=":
		if preferredTemp > *maxTemp {
			return -1
		}
		if preferredTemp > currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp > *minTemp {
			*minTemp = preferredTemp
		}
	case "<=":
		if preferredTemp < *minTemp {
			return -1
		}
		if preferredTemp < currTemp {
			currTemp = preferredTemp
		}
		if preferredTemp < *maxTemp {
			*maxTemp = preferredTemp
		}
	default:
		return -1
	}

	return currTemp
}
