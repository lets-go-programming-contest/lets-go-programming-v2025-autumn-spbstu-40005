package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
	decimalBase    = 10
	minPartsLength = 2
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	departments := parseInt(scanner.Text())

	for range make([]struct{}, departments) {
		scanner.Scan()
		workers := parseInt(scanner.Text())

		minTemp := InitialMinTemp
		maxTemp := InitialMaxTemp

		for range make([]struct{}, workers) {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)

			if len(parts) < minPartsLength {
				fmt.Println(-1)

				continue
			}

			operator := parts[0]
			temp := parseInt(parts[1])

			switch operator {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}

func parseInt(s string) int {
	num := 0
	bytes := []byte(s)

	for _, ch := range bytes {
		if ch >= '0' && ch <= '9' {
			num = num*decimalBase + int(ch-'0')
		}
	}

	return num
}
