package main

import (
	"container/heap"
	"errors"
	"fmt"

	"anton.mezentsev/task-2-2/internal/intheap"
)

var (
	ErrInvalidInput  = errors.New("invalid input parameters")
	ErrEmptyHeap     = errors.New("heap is empty after processing")
	ErrTypeAssertion = errors.New("type assertion failed")
)

func FindKthPreference(ratings []int, preferenceOrder int) (int, error) {
	if len(ratings) == 0 || preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		return 0, ErrInvalidInput
	}

	heapContainer := &intheap.CustomHeap{}
	heap.Init(heapContainer)

	for _, rating := range ratings {
		if heapContainer.Len() < preferenceOrder {
			heap.Push(heapContainer, rating)
		} else if rating > (*heapContainer)[0] {
			(*heapContainer)[0] = rating
			heap.Fix(heapContainer, 0)
		}
	}

	if heapContainer.Len() == 0 {
		return 0, ErrEmptyHeap
	}

	result, ok := heap.Pop(heapContainer).(int)
	if !ok {
		return 0, ErrTypeAssertion
	}

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
	if err != nil {
		fmt.Printf("Invalid preference order: %v\n", err)

		return
	}

	if selectionIndex <= 0 || selectionIndex > totalItems {
		fmt.Printf("Invalid preference order: must be between 1 and %d, got %d\n", totalItems, selectionIndex)

		return
	}

	finalChoice, err := FindKthPreference(scores, selectionIndex)
	if err != nil {
		fmt.Printf("Error finding preference: %v\n", err)

		return
	}

	fmt.Println(finalChoice)
}
