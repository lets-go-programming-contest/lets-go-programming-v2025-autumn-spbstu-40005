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
	length := len(oldSlice)
	if length == 0 {
		return nil
	}
	lastElement := oldSlice[length-1]
	*h = oldSlice[0 : length-1]
	return lastElement
}

func FindKthPreference(ratings []int, preferenceOrder int) int {
	if len(ratings) == 0 || preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		return -1
	}

	heapContainer := &CustomHeap{}
	heap.Init(heapContainer)
	for _, currentRating := range ratings {
		if heapContainer.Len() < preferenceOrder {
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
