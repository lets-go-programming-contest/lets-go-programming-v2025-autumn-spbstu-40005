package queue

type Element struct {
	Value    int
	Priority int
	Index    int
}

type Queue []*Element

func (x Queue) Len() int {
	return len(x)
}

func (x Queue) Less(i, j int) bool {
	return x[i].Priority > x[j].Priority
}

func (x Queue) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
	x[i].Index = i
	x[j].Index = j
}

func (x *Queue) Push(newElement interface{}) {
	n := len(*x)
	item := newElement.(*Element)
	item.Index = n
	*x = append(*x, item)
}

func (x *Queue) Pop() interface{} {
	old := *x
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*x = old[0 : n-1]
	return item
}
