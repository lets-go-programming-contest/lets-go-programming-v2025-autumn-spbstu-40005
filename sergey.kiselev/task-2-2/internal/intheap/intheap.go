package intheap

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	val, ok := x.(int)
	if !ok {
		panic("value is not an int")
	}

	*h = append(*h, val)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
