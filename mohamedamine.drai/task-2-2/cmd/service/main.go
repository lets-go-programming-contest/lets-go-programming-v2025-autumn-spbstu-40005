package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrReadInput       = errors.New("failed to read input")
	ErrPreferenceRange = errors.New("preference out of range")
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
	// Since we control all usage and only push integers,
	// we can safely assume the type is correct
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)

	if n == 0 {
		return nil
	}

	val := old[n-1]
	*h = old[:n-1]

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

	result := getPreference(dishRatings, dishPreference)
	fmt.Println(result)
}

func getPreference(dishRatings []int, preference int) int {
	ratingHeap := &IntHeap{}
	heap.Init(ratingHeap)

	for _, rating := range dishRatings {
		heap.Push(ratingHeap, rating)

		if ratingHeap.Len() > preference {
			heap.Pop(ratingHeap)
		}
	}

	return heap.Pop(ratingHeap).(int)
}
