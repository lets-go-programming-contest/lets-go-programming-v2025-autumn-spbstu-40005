package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrReadInput       = errors.New("failed to read input")
	ErrPreferenceRange = errors.New("preference out of range")
	ErrHeapOperation   = errors.New("heap operation failed")
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	length := len(old)

	if length == 0 {
		return nil
	}

	val := old[length-1]
	*h = old[:length-1]
	return val
}

func main() {
	var (
		dishCount      int
		dishPreference int
	)

	if _, err := fmt.Scan(&dishCount); err != nil {
		fmt.Println("Failed to read dish count:", err)
		return
	}

	dishRatings := make([]int, dishCount)
	for index := range dishRatings {
		if _, err := fmt.Scan(&dishRatings[index]); err != nil {
			fmt.Println("Failed to read dish rating:", err)
			return
		}
	}

	if _, err := fmt.Scan(&dishPreference); err != nil {
		fmt.Println("Failed to read dish preference:", err)
		return
	}

	if dishPreference < 1 || dishPreference > dishCount {
		fmt.Println("Preference out of range")
		return
	}

	result, err := getPreference(dishRatings, dishPreference)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}

func getPreference(dishRatings []int, preference int) (int, error) {
	ratingHeap := &IntHeap{}
	heap.Init(ratingHeap)

	for _, rating := range dishRatings {
		heap.Push(ratingHeap, rating)

		if ratingHeap.Len() > preference {
			heap.Pop(ratingHeap)
		}
	}

	result := heap.Pop(ratingHeap)
	if result == nil {
		return 0, ErrHeapOperation
	}

	value, ok := result.(int)
	if !ok {
		return 0, ErrHeapOperation
	}

	return value, nil
}
