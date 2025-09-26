package main

import (
	"container/heap"

	"sergey.kiselev/task-2-2/internal/maxheap"
)

func findLargest(nums []int, k int) int {
	dishesHeap := &maxheap.MaxHeap{}
	heap.Init(dishesHeap)

	for _, num := range nums {
		heap.Push(dishesHeap, num)
	}

	for range k - 1 {
		heap.Pop(dishesHeap)
	}

	return (*dishesHeap)[0]
}

func main() {}
