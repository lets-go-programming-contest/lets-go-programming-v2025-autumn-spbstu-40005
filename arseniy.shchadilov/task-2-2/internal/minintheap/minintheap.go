package minintheap

import "container/heap"

type MinIntHeap []int

func (h MinIntHeap) Len() int { 
	return len(h) 
}
func (h MinIntHeap) Less(i, j int) bool { 
	return h[i] < h[j] 
}
func (h MinIntHeap) Swap(i, j int) { 
	h[i], h[j] = h[j], h[i] 
}

func (h *MinIntHeap) Push(x interface{}) { 
	*h = append(*h, x.(int)) 
}
func (h *MinIntHeap) Pop() interface{} {
    old := *h
    x := old[len(old)-1]
    *h = old[:len(old)-1]
    return x
}

func KthLargest(values []int, k int) int {
    heap := &MinIntHeap{}
    heap.Init(heap)
    
    for _, value := range values {
        if heap.Len() < k {
            heap.Push(heap, value)
        } else if value > (*heap)[0] {
            heap.Pop(heap)
            heap.Push(heap, value)
        }
    }
    return (*heap)[0]
}