package main

import "container/heap"
import "errors"
import "fmt"

var ErrHeapSize = errors.New("heap size is smaller than the required number")

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
	*heap = append(*heap, x.(int))
}

func (heap *DishHeap) Pop() any {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]
	return x
}

func findKthLargest(dishHeap *DishHeap, k int) (int, error) {
	if dishHeap.Len() < k {
		return 0, ErrHeapSize
	}

	heapCopy := make(DishHeap, dishHeap.Len())
	copy(heapCopy, *dishHeap)
	heap.Init(&heapCopy)

	for i := 1; i < k; i++ {
		heap.Pop(&heapCopy)
	}

	kthLargest := heap.Pop(&heapCopy)
	if value, ok := kthLargest.(int); ok {
		return value, nil
	}

	return 0, errors.New("Failed to retrieve value from heap")
}

func main() {
	var numberOfDishes, k int

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

	if _, err := fmt.Scan(&k); err != nil {
		fmt.Println("Error reading k:", err)
		return
	}

	if k < 1 || k > numberOfDishes {
		fmt.Println("Invalid value for k")
		return
	}

	result, err := findKthLargest(dishHeap, k)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(result)
}
