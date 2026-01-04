package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishCount, kValue int

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		return
	}

	ratings := make([]int, dishCount)
	for index := range ratings {
		_, err = fmt.Scan(&ratings[index])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return
	}

	if kValue < 1 || kValue > dishCount {
		return
	}

	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for _, rating := range ratings[:kValue] {
		heap.Push(minHeap, rating)
	}

	for _, rating := range ratings[kValue:] {
		if rating > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, rating)
		}
	}

	result := (*minHeap)[0]
	fmt.Println(result)
}
