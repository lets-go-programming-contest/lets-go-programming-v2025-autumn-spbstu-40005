package intheap

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(i, j int) bool {
	if (i < 0) || (i >= heap.Len()) {
		panic("first index out of range")
	} else if (j < 0) || (j >= heap.Len()) {
		panic("second index out of range")
	}

	return (*heap)[i] > (*heap)[j]
}

func (heap *IntHeap) Swap(i, j int) {
	if (i < 0) || (i >= heap.Len()) {
		panic("first index out of range")
	} else if (j < 0) || (j >= heap.Len()) {
		panic("second index out of range")
	}

	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *IntHeap) Push(elem any) {
	val, ok := elem.(int)
	if !ok {
		panic("pushed value is not int")
	}

	*heap = append(*heap, val)
}

func (heap *IntHeap) Pop() any {
	if heap.Len() == 0 {
		return nil
	}

	elem := (*heap)[len(*heap)-1]
	*heap = (*heap)[:len(*heap)-1]

	return elem
}
