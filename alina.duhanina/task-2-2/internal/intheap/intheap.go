package intheap

type IntHeap []int

func (heapInstance *IntHeap) Len() int {
	return len(*heapInstance)
}

func (heapInstance *IntHeap) Less(i, j int) bool {
	return (*heapInstance)[i] < (*heapInstance)[j]
}

func (heapInstance *IntHeap) Swap(i, j int) {
	(*heapInstance)[i], (*heapInstance)[j] = (*heapInstance)[j], (*heapInstance)[i]
}

func (heapInstance *IntHeap) Push(x interface{}) {
	num, ok := x.(int)
	if !ok {
		return
	}

	*heapInstance = append(*heapInstance, num)
}

func (heapInstance *IntHeap) Pop() interface{} {
	if len(*heapInstance) == 0 {
		return nil
	}

	old := *heapInstance
	n := len(old)
	x := old[n-1]
	*heapInstance = old[0 : n-1]

	return x
}
