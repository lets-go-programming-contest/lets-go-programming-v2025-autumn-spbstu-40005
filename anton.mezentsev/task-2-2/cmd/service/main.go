package main

import (
	"container/heap"
	"errors"
	"fmt"

	"anton.mezentsev/task-2-2/internal/intheap"
)

func FindKthPreference(ratings []int, preferenceOrder int) (int, error) {
	if len(ratings) == 0 || preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		return 0, errors.New("invalid input parameters")
	}

	heapContainer := &intheap.CustomHeap{}
	heap.Init(heapContainer)

	for _, rating := range ratings {
		if heapContainer.Len() < preferenceOrder {
			heap.Push(heapContainer, rating)
		} else if rating > (*heapContainer)[0] {
			heap.Pop(heapContainer)
			heap.Push(heapContainer, rating)
		}
	}

	if heapContainer.Len() == 0 {
		return 0, errors.New("heap is empty after processing")
	}

	result := heap.Pop(heapContainer).(int)
	return result, nil
}

func main() {
	var totalItems int

	_, err := fmt.Scan(&totalItems)
	if err != nil {
		fmt.Printf("Error reading number of dishes: %v\n", err)

		return
	}
	if totalItems <= 0 {
		fmt.Printf("Invalid number of dishes: must be positive, got %d\n", totalItems)

		return
	}

	scores := make([]int, totalItems)
	for index := range totalItems {
		_, err := fmt.Scan(&scores[index])
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

	finalChoice, err := FindKthPreference(scores, selectionIndex)
	if err != nil {
		fmt.Printf("Error finding preference: %v\n", err)
		return
	}
	fmt.Println(finalChoice)
}
