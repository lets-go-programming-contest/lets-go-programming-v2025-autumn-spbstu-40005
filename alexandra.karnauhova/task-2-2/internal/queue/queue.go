package queue

type Element struct {
	Value    int
	Priority int
	Index    int
}

type Queue []*Element

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Less(i, j int) bool {
	return (*q)[i].Priority > (*q)[j].Priority
}

func (q *Queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
	(*q)[i].Index = i
	(*q)[j].Index = j
}

func (q *Queue) Push(newElement interface{}) {
	element, oke := newElement.(*Element)
	if !oke {
		return
	}
	n := len(*q)
	element.Index = n
	*q = append(*q, element)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*q = old[0 : n-1]

	return item
}
