package main

import "container/heap"

type MinHeap []int

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(i, j int) bool {
	if i < 0 || i >= len(*h) || j < 0 || j >= len(*h) {
		panic("index out of ranged")
	}

	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Swap(i, j int) {
	if i < 0 || i >= len(*h) || j < 0 || j >= len(*h) {
		panic("index out of range")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("push error")
	}

	*h = append(*h, value)
}

func (h *MinHeap) Pop() interface{} {
	if len(*h) == 0 {
		panic("pop error")
	}

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
		panic("popvalue error")
	}

	return intValue
}

func (h *MinHeap) Peek() int {
	return (*h)[0]
}
