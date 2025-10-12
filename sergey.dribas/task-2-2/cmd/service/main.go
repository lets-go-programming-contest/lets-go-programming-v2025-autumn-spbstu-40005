package main

import (
	"container/heap"
	"errors"
	"fmt"

	"sergey.dribas/task-2-2/internal/intheap"
)

var (
	ErrHeap = errors.New("heap size less than number")
	ErrType = errors.New("unexpected type from heap pop")
)

func FindKthSmallest(intHeap *intheap.IntHeap, number int) (int, error) {
	if intHeap.Len() < number {
		return 0, ErrHeap
	}

	copyHeap := make(intheap.IntHeap, intHeap.Len())
	copy(copyHeap, *intHeap)

	heap.Init(&copyHeap)

	for range copyHeap.Len() - number {
		heap.Pop(&copyHeap)
	}

	result := heap.Pop(&copyHeap)
	if val, ok := result.(int); ok {
		return val, nil
	}

	return 0, ErrType
}

func main() {
	var (
		size int
		dish = &intheap.IntHeap{}
	)

	heap.Init(dish)

	if _, err := fmt.Scan(&size); err != nil {
		return
	}

	var number int
	for range size {
		if _, err := fmt.Scan(&number); err != nil {
			return
		}

		dish.Push(number)
	}

	var predict int
	if _, err := fmt.Scan(&predict); err != nil {
		return
	}

	if result, err := FindKthSmallest(dish, predict); err == nil {
		fmt.Println(result)
	}
}
