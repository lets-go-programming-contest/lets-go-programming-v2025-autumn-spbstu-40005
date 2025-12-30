package intheap

import (
	"container/heap"
	"fmt"
)

type CustomHeap []int

func (h *CustomHeap) Len() int {
	return len(*h)
}

func (h *CustomHeap) Less(i, j int) bool {
	if i >= h.Len() || j >= h.Len() {
		panic(fmt.Sprintf("index out of range: i=%d, j=%d, len=%d", i, j, h.Len()))
	}

	return (*h)[i] < (*h)[j]
}

func (h *CustomHeap) Swap(i, j int) {
	if i >= h.Len() || j >= h.Len() {
		panic(fmt.Sprintf("index out of range: i=%d, j=%d, len=%d", i, j, h.Len()))
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *CustomHeap) Push(x interface{}) {
	value, correctType := x.(int)
	if !correctType {
		panic(fmt.Sprintf("invalid type: expected int, got %T", x))
	}

	*h = append(*h, value)
}

func (h *CustomHeap) Pop() interface{} {
	if h.Len() == 0 {
		panic("pop from empty heap")
	}

	oldSlice := *h
	n := len(oldSlice)
	lastElement := oldSlice[n-1]
	*h = oldSlice[0 : n-1]

	return lastElement
}

func FindKthPreference(ratings []int, preferenceOrder int) int {
	if len(ratings) == 0 || preferenceOrder <= 0 || preferenceOrder > len(ratings) {
		panic(fmt.Sprintf("invalid parameters: ratings len=%d, preference order=%d", len(ratings), preferenceOrder))
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

	return (*heapContainer)[0]
}
