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

	parts := strings.Fields(preferences[0])
	preferedTemp, _ := strconv.Atoi(parts[1])
	sign := parts[0]
	if preferedTemp < MinTemp || preferedTemp > MaxTemp {
		currTemp = -1
	} else if sign == ">=" {
		currTemp = preferedTemp
	} else {
		currTemp = MinTemp
	}

	for idx := 1; idx < len(preferences); idx++ {
		parts := strings.Fields(preferences[idx])
		preferedTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]

		if preferedTemp < MinTemp || preferedTemp > MaxTemp {
			currTemp = -1

		} else {
			switch sign {
			case ">=":
				if preferedTemp > maxTemp {
					currTemp = -1
				} else if preferedTemp > currTemp {
					currTemp = preferedTemp
				}

				if preferedTemp > minTemp {
					minTemp = preferedTemp
				}
			case "<=":
				if preferedTemp < minTemp {
					currTemp = -1
				} else if preferedTemp < currTemp {
					currTemp = preferedTemp
				}

				if preferedTemp < maxTemp {
					maxTemp = preferedTemp
				}
			default:
				currTemp = -1
			}
		}

		if currTemp == -1 {
			for range len(preferences) - idx {
				fmt.Println(currTemp)
			}

			return
		}

		fmt.Println(currTemp)
	}
}
