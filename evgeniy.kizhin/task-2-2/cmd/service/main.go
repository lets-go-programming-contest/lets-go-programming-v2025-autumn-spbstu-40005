package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrHeapSize         = errors.New("heap size is smaller than the required number")
	ErrFailedToGetValue = errors.New("failed to get value from heap")
)

type DishHeap []int

func (heap *DishHeap) Len() int {
	return len(*heap)
}

func (heap *DishHeap) Less(i, j int) bool {
	return (*heap)[i] > (*heap)[j]
}

func (heap *DishHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *DishHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("Invalid type")
	}

	*heap = append(*heap, value)
}

func (heap *DishHeap) Pop() any {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

func findKthLargest(dishHeap *DishHeap, kthPos int) (int, error) {
	if dishHeap.Len() < kthPos {
		return 0, ErrHeapSize
	}

	heapCopy := make(DishHeap, dishHeap.Len())
	copy(heapCopy, *dishHeap)
	heap.Init(&heapCopy)

	for i := 1; i < kthPos; i++ {
		heap.Pop(&heapCopy)
	}

	kthLargest := heap.Pop(&heapCopy)
	if value, ok := kthLargest.(int); ok {
		return value, nil
	}

	return 0, ErrFailedToGetValue
}

func main() {
	var numberOfDishes, kthPos int

	if _, err := fmt.Scan(&numberOfDishes); err != nil {
		fmt.Println("Error reading number of dishes:", err)

		return
	}

	dishHeap := &DishHeap{}
	heap.Init(dishHeap)

	for range numberOfDishes {
		var dishPreference int
		if _, err := fmt.Scan(&dishPreference); err != nil {
			fmt.Println("Error reading dish preference:", err)

			return
		}

		heap.Push(dishHeap, dishPreference)
	}

	if _, err := fmt.Scan(&kthPos); err != nil {
		fmt.Println("Error reading k:", err)

		return
	}

	if kthPos < 1 || kthPos > numberOfDishes {
		fmt.Println("Invalid value for k")

		return
	}

	result, err := findKthLargest(dishHeap, kthPos)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(result)
}
