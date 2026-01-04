package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}

	firstLine := scanner.Text()
	parts := strings.Fields(firstLine)
	if len(parts) < 2 {
		return
	}

	N, _ := strconv.Atoi(parts[0])
	K, _ := strconv.Atoi(parts[1])

	for i := 0; i < N; i++ {
		minTemp := 15
		maxTemp := 30

		for j := 0; j < K; j++ {
			if !scanner.Scan() {
				return
			}

			request := scanner.Text()
			requestParts := strings.Fields(request)
			if len(requestParts) < 2 {
				continue
			}

			operator := requestParts[0]
			temp, _ := strconv.Atoi(requestParts[1])

			if operator == ">=" {
				if temp > minTemp {
					minTemp = temp
				}
			} else if operator == "<=" {
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(maxTemp)
			}
		}
	}
}
