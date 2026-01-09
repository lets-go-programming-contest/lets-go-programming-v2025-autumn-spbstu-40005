package dishheap

type DishHeap []int

func (h *DishHeap) Len() int {
	return len(*h)
}

func (h *DishHeap) Less(i int, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *DishHeap) Swap(i int, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *DishHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("expected int")
	}

	*h = append(*h, value)
}

func (h *DishHeap) Pop() interface{} {
	if len(*h) == 0 {
		return nil
	}

	last := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]

	return last
}
