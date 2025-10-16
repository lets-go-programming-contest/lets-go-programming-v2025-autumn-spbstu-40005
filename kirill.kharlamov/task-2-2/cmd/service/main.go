package main

import (
	"fmt"

	"kirill.kharlamov/task-2-2/internal/intheap"
)

func main() {
	var totalItems int

	_, err := fmt.Scan(&totalItems)
	if err != nil || totalItems <= 0 {
		fmt.Printf("Invalid number of dishes: %v\n", err)
		return
	}

	scores := make([]int, totalItems)
	for i := range totalItems {
		_, err := fmt.Scan(&scores[i])
		if err != nil {
			fmt.Printf("Invalid rating: %v\n", err)
			return
		}
	}

	var selectionIndex int
	_, err = fmt.Scan(&selectionIndex)
	if err != nil || selectionIndex <= 0 || selectionIndex > totalItems {
		fmt.Printf("Invalid preference order: %v\n", err)
		return
	}

	finalChoice := intheap.FindKthPreference(scores, selectionIndex)
	fmt.Println(finalChoice)
}
