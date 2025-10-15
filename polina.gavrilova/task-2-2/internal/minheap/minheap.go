package minheap

import "container/heap"

type MinHeap []int

func (heap *MinHeap) Len() int {
	return len(*heap)
}

func (heap *MinHeap) Less(i, j int) bool {
	return (*heap)[i] < (*heap)[j]
}

func (heap *MinHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*heap = append(*heap, value)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func FindKthLargest(ratings []int, preferenceOrder int) int {
	heapInstance := &MinHeap{}
	heap.Init(heapInstance)

	for _, rating := range ratings {
		if heapInstance.Len() < preferenceOrder {
			heap.Push(heapInstance, rating)
		} else if rating > (*heapInstance)[0] {
			heap.Pop(heapInstance)
			heap.Push(heapInstance, rating)
		}
	}

	return (*heapInstance)[0]
}
