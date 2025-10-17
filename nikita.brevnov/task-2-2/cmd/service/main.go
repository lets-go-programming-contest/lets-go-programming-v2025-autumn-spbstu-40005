package main

import (
	"container/heap"
	"errors"

	"nikita.brevnov/task-2-2/internal/intheap"
)

var (
	errEmptyHeap = errors.New("the heap is already empty")
	errConvert   = errors.New("not integer in heap")
)

func findLargest(nums []int, number int) (int, error) {
	dishesHeap := &intheap.IntHeap{}
	heap.Init(dishesHeap)

	for _, num := range nums {
		heap.Push(dishesHeap, num)
	}

	for range number - 1 {
		heap.Pop(dishesHeap)
	}

	val := heap.Pop(dishesHeap)
	if val == nil {
		return 0, errEmptyHeap
	}

	value, ok := val.(int)
	if !ok {
		return 0, errConvert
	}

	return value, nil
}

func main() {
}
