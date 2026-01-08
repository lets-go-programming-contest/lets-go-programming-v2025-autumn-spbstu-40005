package main

import (
	"fmt"

	"github.com/dmitei/task-2-2/internal/heap"
)

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Println("scan dishesCount:", err)
		return
	}

	ratings := make([]int, dishesCount)
	for i := range ratings {
		if _, err := fmt.Scan(&ratings[i]); err != nil {
			fmt.Println("scan rating:", err)
			return
		}
	}

	var position int
	if _, err := fmt.Scan(&position); err != nil {
		fmt.Println("scan position:", err)
		return
	}

	result, err := heap.FindKthPreferred(ratings, position)
	if err != nil {
		fmt.Println("find kth preferred:", err)
		return
	}

	fmt.Println(result)
}
