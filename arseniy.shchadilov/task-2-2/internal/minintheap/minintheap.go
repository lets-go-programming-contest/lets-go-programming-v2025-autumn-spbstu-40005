package minintheap

import "container/heap"

type MinIntHeap []int

func (h *MinIntHeap) Len() int {
	return len(*h)
}

func (h *MinIntHeap) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= len(*h) || j >= len(*h) {
		panic("MinIntHeap.Less: index out of range")
	}

	return (*h)[i] < (*h)[j]
}

func (h *MinIntHeap) Swap(i, j int) {
	if i < 0 || j < 0 || i >= len(*h) || j >= len(*h) {
		panic("MinIntHeap.Swap: index out of range")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinIntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("MinIntHeap.Push: expected int")
	}

	*h = append(*h, value)
}

func (h *MinIntHeap) Pop() interface{} {
	if len(*h) == 0 {
		panic("MinIntHeap.Pop: heap is empty")
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func KthLargest(values []int, kth int) int {
	if len(values) == 0 || kth <= 0 || kth > len(values) {
		panic("KthLargest: invalid arguments - empty values or kth out of range")
	}

	heapData := &MinIntHeap{}

	for _, value := range values {
		if heapData.Len() < kth {
			heap.Push(heapData, value)
		} else if value > (*heapData)[0] {
			heap.Pop(heapData)
			heap.Push(heapData, value)
		}
	}

	if heapData.Len() == 0 {
		panic("KthLargest: unexpected empty heap")
	}

	return (*heapData)[0]
}
