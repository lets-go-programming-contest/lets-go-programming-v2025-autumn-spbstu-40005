package intheap

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	if h.Len() < i+1 || h.Len() < j+1 || i < 0 || j < 0 {
		panic("out of range")
	}

	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	if h.Len() < i+1 || h.Len() < j+1 || j < 0 || i < 0 {
		panic("out of range")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	val, ok := x.(int)
	if ok {
		*h = append(*h, val)
	} else {
		panic("push value isn`t int")
	}
}

func (h *IntHeap) Pop() interface{} {
	if h.Len() == 0 {
		return nil
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
