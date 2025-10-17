package main

import (
	"container/heap"
	"fmt"

	"feodor.khoroshilov/task-2-2/internal/intheap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	var totalDishesCount int
	if _, err := fmt.Scan(&totalDishesCount); err != nil {
		fmt.Printf("Error reading total dishes count: %v\n", err)

		return
	}

	if totalDishesCount < 1 || totalDishesCount > 10000 {
		fmt.Println("Error: number of dishes must be between 1 and 10000")

		return
	}

	dishRatingsHeap := &intheap.IntHeap{}
	heap.Init(dishRatingsHeap)

	for i := 0; i < totalDishesCount; i++ {
		var dishRating int
		if _, err := fmt.Scan(&dishRating); err != nil {
			fmt.Printf("Error reading dish rating: %v\n", err)

			return
		}

		if dishRating < -10000 || dishRating > 10000 {
			fmt.Println("Error: dish rating must be between -10000 and 10000")

			return
		}

		heap.Push(dishRatingsHeap, dishRating)
	}

	var preferredDishPosition int
	if _, err := fmt.Scan(&preferredDishPosition); err != nil {
		fmt.Printf("Error reading preferred dish position: %v\n", err)

		return
	}

	if preferredDishPosition < 1 || preferredDishPosition > totalDishesCount {
		fmt.Printf("Error: preferred dish position must be between 1 and %d\n", totalDishesCount)
		
		return
	}

	for range preferredDishPosition-1 {
		heap.Pop(dishRatingsHeap)
	}

	preferredDishRating := heap.Pop(dishRatingsHeap)
	fmt.Println(preferredDishRating)
}