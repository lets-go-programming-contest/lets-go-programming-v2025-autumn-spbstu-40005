package heap

import "fmt"

var errPositionOutOfRange = fmt.Errorf("position out of range")

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
	oldSlice := *h
	length := len(oldSlice)
	if length == 0 {
		return nil
	}

	value := oldSlice[length-1]
	*h = oldSlice[:length-1]

	return value
}

func FindKthPreferred(ratings []int, position int) (int, error) {
	if position <= 0 || position > len(ratings) {
		return 0, fmt.Errorf("%w: %d not in [1,%d]", errPositionOutOfRange, position, len(ratings))
	}

	heapInstance := &MinHeap{}
	for _, rating := range ratings {
		heapInstance.Push(rating)

		if heapInstance.Len() > position {
			_ = heapInstance.Pop()
		}
	}

	result := (*heapInstance)[0]

	return result, nil
}
