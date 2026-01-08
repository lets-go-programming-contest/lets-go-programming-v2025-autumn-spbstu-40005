package heap

import "fmt"

type MinHeap []int

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic(fmt.Sprintf("expected int, got %T", x))
	}

	*h = append(*h, value)
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	if n == 0 {
		return nil
	}

	value := old[n-1]
	*h = old[:n-1]

	return value
}

func FindKthPreferred(ratings []int, position int) (int, error) {
	if position <= 0 || position > len(ratings) {
		return 0, fmt.Errorf("position %d out of range 1..%d", position, len(ratings))
	}

	h := &MinHeap{}
	for _, r := range ratings {
		h.Push(r)

		if h.Len() > position {
			_ = h.Pop()
		}
	}

	result := (*h)[0]

	return result, nil
}
