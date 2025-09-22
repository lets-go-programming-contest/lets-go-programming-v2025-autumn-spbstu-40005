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
	for _ = range number {
		_, err = fmt.Scan(&count)
		if err != nil {
			return
		}
		var (
			sign         string
			preferedTemp int
			preferences  []string
		)
		for _ = range count {
			_, err = fmt.Scan(&sign, &preferedTemp)
			if err != nil {
				return
			}
			var temp string = sign + " " + strconv.Itoa(preferedTemp)
			preferences = append(preferences, temp)
		}
		changeTemperature(preferences)
	}
}

func changeTemperature(preferences []string) {
	const (
		MAX_TEMP = 30
		MIN_TEMP = 0
	)
	var (
		maxTemp  int = MAX_TEMP
		minTemp  int = MIN_TEMP
		currTemp int
	)
	for i := 0; i < len(preferences); i++ {
		parts := strings.Fields(preferences[i])
		preferedTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]
		if preferedTemp < 15 || preferedTemp > 30 {
			currTemp = -1
		} else {
			if i == 0 {
				currTemp = preferedTemp
			}
			if sign == ">=" {
				if preferedTemp > maxTemp {
					currTemp = -1
				} else if preferedTemp > currTemp {
					currTemp = preferedTemp
				}
				if preferedTemp > minTemp {
					minTemp = preferedTemp
				}
			} else if sign == "<=" {
				if preferedTemp < minTemp {
					currTemp = -1
				} else if preferedTemp < currTemp {
					currTemp = preferedTemp
				}
				if preferedTemp < maxTemp {
					maxTemp = preferedTemp
				}
			} else {
				currTemp = -1
			}
		}
		fmt.Println(currTemp)
	}
}
