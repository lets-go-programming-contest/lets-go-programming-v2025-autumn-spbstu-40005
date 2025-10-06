package intheap

import (
	"container/heap"
)

type IntHeap []int

func (h *IntHeap) Push(val any) {
	intValue, err := val.(int)
	if !err {
		panic("invalid type pushed to IntHeap")
	}

	*h = append(*h, intValue)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	orig := *h
	origLength := len(orig)
	toreturn := orig[origLength-1]
	*h = orig[0 : origLength-1]

	return toreturn
}

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) GetNth(numberOfElement int) int {
	temp := make(IntHeap, h.Len())
	copy(temp, *h)
	heap.Init(&temp)

	for range numberOfElement - 1 {
		heap.Pop(&temp)
	}

	return heap.Pop(&temp).(int)
}
