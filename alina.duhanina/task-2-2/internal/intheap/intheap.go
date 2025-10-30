package intheap

type IntHeap []int

func (heapInstance *IntHeap) Len() int {
	return len(*heapInstance)
}

func (heapInstance *IntHeap) Less(i, j int) bool {
	if i < 0 || i >= heapInstance.Len() || j < 0 || j >= heapInstance.Len() {
		panic("index out of range")
	}

	return (*heapInstance)[i] > (*heapInstance)[j]
}

func (heapInstance *IntHeap) Swap(i, j int) {
	if i < 0 || i >= heapInstance.Len() || j < 0 || j >= heapInstance.Len() {
		panic("index out of range")
	}

	(*heapInstance)[i], (*heapInstance)[j] = (*heapInstance)[j], (*heapInstance)[i]
}

func (heapInstance *IntHeap) Push(x interface{}) {
	num, ok := x.(int)
	if !ok {
		panic("invalid type for heap")
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
