package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	MinTemp = 15
	MaxTemp = 30
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	N, _ := strconv.Atoi(scanner.Text())

	if !scanner.Scan() {
		return
	}
	K, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < N; i++ {
		minTemp := MinTemp
		maxTemp := MaxTemp

		for j := 0; j < K; j++ {
			if !scanner.Scan() {
				return
			}
			operator := scanner.Text()

			if !scanner.Scan() {
				return
			}
			value, _ := strconv.Atoi(scanner.Text())

			switch operator {
			case ">=":
				if value > minTemp {
					minTemp = value
				}
			case "<=":
				if value < maxTemp {
					maxTemp = value
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
