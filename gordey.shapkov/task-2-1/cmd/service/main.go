package main

import (
	"fmt"
)

type TemperaturePreference struct {
	maxTemp, minTemp, currTemp int
}

const (
	MaxTemp = 30
	MinTemp = 15
)

func main() {
	var number, count int
	if _, err := fmt.Scan(&number); err != nil {
		return
	}

	for range number {
		_, err := fmt.Scan(&count)
		if err != nil {
			return
		}

		var (
			sign         string
			preferedTemp int
			temp         TemperaturePreference
		)

		temp.maxTemp = MaxTemp
		temp.minTemp = MinTemp
		temp.currTemp = MinTemp

		for range count {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}

			if temp.currTemp == -1 {
				fmt.Println(temp.currTemp)

				continue
			}

			changeTemperature(sign, preferedTemp, &temp)
			fmt.Println(temp.currTemp)
		}
	}
}

func changeTemperature(sign string, preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp < MinTemp || preferedTemp > MaxTemp {
		temp.currTemp = -1
	}

	switch sign {
	case ">=":
		handleGreaterEqual(preferedTemp, temp)
	case "<=":
		handleLessEqual(preferedTemp, temp)
	default:
		temp.currTemp = -1
	}
}

func handleGreaterEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp > temp.maxTemp {
		temp.currTemp = -1

		return
	}

	if preferedTemp > temp.currTemp {
		temp.currTemp = preferedTemp
	}

	if preferedTemp > temp.minTemp {
		temp.minTemp = preferedTemp
	}
}

func handleLessEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp < temp.minTemp {
		temp.currTemp = -1

		return
	}

	if preferedTemp < temp.currTemp {
		temp.currTemp = preferedTemp
	}

	if preferedTemp < temp.maxTemp {
		temp.maxTemp = preferedTemp
	}
}
