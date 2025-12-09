package main

import (
	"container/heap"
	"errors"
	"fmt"

	"polina.gavrilova/task-2-2/internal/minheap"
)

var ErrPreferenceOrder = errors.New("invalid preference order")

func FindKthLargest(ratings []int, preferenceOrder int) (int, error) {
	if preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		return 0, ErrPreferenceOrder
	}

	heapInstance := &minheap.MinHeap{}
	heap.Init(heapInstance)

	for _, rating := range ratings {
		if heapInstance.Len() < preferenceOrder {
			heap.Push(heapInstance, rating)
		} else {
			top, err := heapInstance.Top()
			if err != nil {
				return 0, fmt.Errorf("FindKthLargest: get top during processing: %w", err)
			}

			if rating > top {
				heap.Pop(heapInstance)
				heap.Push(heapInstance, rating)
			}
		}
	}

	result, err := heapInstance.Top()
	if err != nil {
		return 0, fmt.Errorf("FindKthLargest: get top during processing: %w", err)
	}

	return result, nil
}

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

	result, err := FindKthLargest(ratings, preferenceOrder)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	fmt.Println(result)
}
