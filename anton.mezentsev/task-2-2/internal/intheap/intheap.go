package intheap

import "container/heap"

type CustomHeap []int

func (h CustomHeap) Len() int {

	return len(h)
}

func (h CustomHeap) Less(index1, index2 int) bool {

	return h[index1] < h[index2]
}

func (h CustomHeap) Swap(index1, index2 int) {
	h[index1], h[index2] = h[index2], h[index1]
}

func (h *CustomHeap) Push(x interface{}) {
	value, correnType := x.(int)
	if !correnType {

		return
	}

	*h = append(*h, value)
}

func (h *CustomHeap) Pop() interface{} {
	oldSlice := *h
	n := len(oldSlice)

	if n == 0 {

		return nil
	}

	lastElement := oldSlice[n-1]
	*h = oldSlice[0 : n-1]

	return lastElement
}

func FindKthPreference(ratings []int, k int) int {
	if len(ratings) == 0 || k <= 0 || k > len(ratings) {

		return -1
	}

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

	if heapContainer.Len() == 0 {

		return -1
	}

	return (*heapContainer)[0]
}
