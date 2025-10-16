package intheap

import "container/heap"

type CustomHeap []int

func (h CustomHeap) Len() int {
	return len(h)
}

func (h CustomHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h CustomHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *CustomHeap) Push(x interface{}) {
	value, correctType := x.(int)
	if !correctType {
		return
	}

	*h = append(*h, value)
}

func (h *CustomHeap) Pop() interface{} {
	oldSlice := *h
	n := len(oldSlice)
	lastElement := oldSlice[n-1]
	*h = oldSlice[0 : n-1]

	return lastElement
}

func FindKthPreference(ratings []int, k int) int {
	heapContainer := &CustomHeap{}
	heap.Init(heapContainer)

	for _, currentRating := range ratings {
		if heapContainer.Len() < k {
			heap.Push(heapContainer, currentRating)
		} else if currentRating > (*heapContainer)[0] {
			heap.Pop(heapContainer)
			heap.Push(heapContainer, currentRating)
		}
	}

	return (*heapContainer)[0]
}
