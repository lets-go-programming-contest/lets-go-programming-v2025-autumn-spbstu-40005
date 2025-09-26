package main

import (
	"container/heap"
	"fmt"
	"gordey.shapkov/task-2-2/internal/IntHeap"
)

func main() {
	h := &IntHeap.IntHeap{3, 2, 1, 5, 6, 6}
	heap.Init(h)
	fmt.Println(heap.Remove(h, 1))
}
