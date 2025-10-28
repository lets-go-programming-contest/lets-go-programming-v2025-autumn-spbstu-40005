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
	intNum, ok := x.(int)
	if !ok {
		panic("invalid input: expected int type")
	}

	*h = append(*h, intNum)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	oldHeap := *h
	n := len(oldHeap)
	res := oldHeap[n-1]
	*h = oldHeap[0 : n-1]

	return res
}
