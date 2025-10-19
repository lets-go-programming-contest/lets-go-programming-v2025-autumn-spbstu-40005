package intheap

import (
	"container/heap"
	"errors"
)

var errEmptyHeap = errors.New("empty heap")

type IntHeap []int //nolint:recvcheck

func (h *IntHeap) Push(val any) {
	intVal, ok := val.(int)
	if !ok {
		panic("invalid type pushed to IntHeap")
	}

	*h = append(*h, intVal)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	n := len(*h)
	val := (*h)[n-1]
	*h = (*h)[:n-1]

	return val
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	if i < 0 || i >= h.Len() || j < 0 || j >= h.Len() {
		panic("out of range from IntHeap.Less")
	}

	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	if i < 0 || i >= h.Len() || j < 0 || j >= h.Len() {
		panic("out of range from IntHeap.Swap")
	}

	h[i], h[j] = h[j], h[i]
}

func (h IntHeap) ReplaceTop(value int) error {
	if h.Len() == 0 {
		return errEmptyHeap
	}

	h[0] = value
	heap.Fix(&h, 0)

	return nil
}

func (h IntHeap) Top() (int, error) {
	if h.Len() == 0 {
		return 0, errEmptyHeap
	}

	return h[0], nil
}
