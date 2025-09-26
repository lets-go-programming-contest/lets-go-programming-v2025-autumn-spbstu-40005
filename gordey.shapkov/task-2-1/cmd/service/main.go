package main

import (
	"fmt"
)

type TemperaturePreference struct {
	maxTemp, minTemp, currTemp int
}

func NewTemperaturePreference(maxTemp, minTemp, currTemp int) *TemperaturePreference {
	return &TemperaturePreference{maxTemp, minTemp, currTemp}
}

func (temp *TemperaturePreference) setMaxTemp(maxTemp int) {
	temp.maxTemp = maxTemp
}

func (temp *TemperaturePreference) setMinTemp(minTemp int) {
	temp.minTemp = minTemp
}

func (temp *TemperaturePreference) setCurrTemp(currTemp int) {
	temp.currTemp = currTemp
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
		)

		temp := NewTemperaturePreference(MaxTemp, MinTemp, MinTemp)

		for range count {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}

			if temp.currTemp == -1 {
				fmt.Println(temp.currTemp)

				continue
			}

			temp.changeTemperature(sign, preferedTemp)
			fmt.Println(temp.currTemp)
		}
	}
}

func (temp *TemperaturePreference) changeTemperature(sign string, preferedTemp int) {
	if preferedTemp < MinTemp || preferedTemp > MaxTemp {
		temp.setCurrTemp(-1)
	}

	switch sign {
	case ">=":
		handleGreaterEqual(preferedTemp, temp)
	case "<=":
		handleLessEqual(preferedTemp, temp)
	default:
		temp.setCurrTemp(-1)
	}
}

func handleGreaterEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp > temp.maxTemp {
		temp.setCurrTemp(-1)

		return
	}

	if preferedTemp > temp.currTemp {
		temp.setCurrTemp(preferedTemp)
	}

	if preferedTemp > temp.minTemp {
		temp.setMinTemp(preferedTemp)
	}
}

func handleLessEqual(preferedTemp int, temp *TemperaturePreference) {
	if preferedTemp < temp.minTemp {
		temp.setCurrTemp(-1)

		return
	}

	if preferedTemp < temp.currTemp {
		temp.setCurrTemp(preferedTemp)
	}

	if preferedTemp < temp.maxTemp {
		temp.setMaxTemp(preferedTemp)
	}
}
