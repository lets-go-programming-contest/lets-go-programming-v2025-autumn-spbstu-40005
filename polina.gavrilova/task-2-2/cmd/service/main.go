package main

import (
	"container/heap"
	"fmt"

	"polina.gavrilova/task-2-2/internal/minheap"
)

func FindKthLargest(ratings []int, preferenceOrder int) int {
	if preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		panic("invalid preferenceOrder")
	}

	heapInstance := &minheap.MinHeap{}
	heap.Init(heapInstance)

	for _, rating := range ratings {
		if heapInstance.Len() < preferenceOrder {
			heap.Push(heapInstance, rating)
		} else if rating > heapInstance.Top() {
			heap.Pop(heapInstance)
			heap.Push(heapInstance, rating)
		}
	}

	return heapInstance.Top()
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

	result := FindKthLargest(ratings, preferenceOrder)
	fmt.Println(result)
}
