package intheap

import (
	"container/heap"
	"errors"
	"fmt"
)

var errOutOfRange = errors.New("no such index in heap")

type IntHeap []int

func (h *IntHeap) Push(val any) {
	intValue, ok := val.(int)
	if !ok {
		panic("invalid type pushed to heap")
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
	if i >= len(*h) || j >= len(*h) || i < 0 || j < 0 {
		panic("comparing indexes out of range")
	}

	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	if i >= len(*h) || j >= len(*h) || i < 0 || j < 0 {
		panic("swapping indexes out of range")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) GetNth(numberOfElement int) (int, error) {
	if numberOfElement > h.Len() || numberOfElement < 0 {
		return 0, fmt.Errorf("impossible to get nth element: %w", errOutOfRange)
	}

	temp := make(IntHeap, h.Len())
	copy(temp, *h)
	heap.Init(&temp)

	for range numberOfElement - 1 {
		heap.Pop(&temp)
	}

	intValue, ok := (heap.Pop(&temp)).(int)
	if !ok {
		panic("invalid type to get from heap")
	}

	return intValue, nil
}
