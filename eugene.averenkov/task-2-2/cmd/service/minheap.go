package main

import "container/heap"

// MinHeap тип для минимальной кучи.
type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, value)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func InitMinHeap() *MinHeap {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	return minHeap
}

func (h *MinHeap) PushValue(x int) {
	heap.Push(h, x)
}

func (h *MinHeap) PopValue() int {
	value := heap.Pop(h)
	intValue, ok := value.(int)
	if !ok {
		return 0
	}

	return intValue
}

func (h *MinHeap) Peek() int {
	return (*h)[0]
}
