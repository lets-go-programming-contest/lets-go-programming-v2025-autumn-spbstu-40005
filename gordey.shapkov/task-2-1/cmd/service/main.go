package main

import (
	"fmt"
	"strings"
	"strconv"
)

func main() {
	var N, K int
	fmt.Scan(&N)
	for i := 0; i < N; i++ {
		fmt.Scan(&K)
		var (sign string
		     preferedTemp int
		     preferences []string)
		for j := 0; j < K; j++ {
			fmt.Scan(&sign, &preferedTemp)
			var temp string = sign + " " + strconv.Itoa(preferedTemp)
			preferences = append(preferences, temp)
		}
		changeTemperature(preferences)
	}
}

func changeTemperature(preferences []string) {
	var (maxTemp int = 100
	    minTemp int = 0
	    currTemp int)
	for i := 0; i < len(preferences); i++ {
		parts := strings.Fields(preferences[i])
		preferedTemp, _ := strconv.Atoi(parts[1])
		sign := parts[0]
		if preferedTemp < 15 || preferedTemp > 30 {
			currTemp = -1;
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
