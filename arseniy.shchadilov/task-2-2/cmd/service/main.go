package main

import (
	"fmt"

	"arseniy.shchadilov/task-2-2/internal/minintheap"
)

func main() {
	var dishCount int

	_, err := fmt.Scan(&dishCount)
	if err != nil || dishCount < 1 || dishCount > 10000 {
		fmt.Println("Error: invalid dish count")
		return
	}

	preferences := make([]int, dishCount)
	for i := range preferences {
		_, err := fmt.Scan(&preferences[i])
		if err != nil || preferences[i] < -10000 || preferences[i] > 10000 {
			fmt.Println("Error: invalid preference value")
			return
		}
	}

	var k int
	_, err = fmt.Scan(&k)
	if err != nil || k < 1 || k > dishCount {
		fmt.Println("Error: invalid preference order")
		return
	}

	result := minintheap.KthLargest(preferences, k)
	fmt.Println(result)
}