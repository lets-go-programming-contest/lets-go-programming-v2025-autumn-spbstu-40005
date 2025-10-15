package main

import "fmt"

func main() {
	var nDishes int

	_, err := fmt.Scan(&nDishes)
	if err != nil || nDishes <= 0 {
		fmt.Printf("Invalid number of dishes: %v\n", err)

		return
	}

	ratings := make([]int, nDishes)
	for i := range nDishes {
		_, err := fmt.Scan(&ratings[i])
		if err != nil {
			fmt.Printf("Invalid rating: %v\n", err)
			return
		}
	}

	var preferenceOrder int
	_, err = fmt.Scan(&preferenceOrder)
	if err != nil || preferenceOrder <= 0 || preferenceOrder > nDishes {
		fmt.Printf("Invalid preference order: %v\n", err)
		return
	}

	result := 0
	fmt.Println(result)
}
