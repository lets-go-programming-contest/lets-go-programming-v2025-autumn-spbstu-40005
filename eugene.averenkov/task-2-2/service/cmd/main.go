package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrInvalidN       = errors.New("invalid number")
	ErrInvalidElement = errors.New("invalid element")
	ErrInvalidK       = errors.New("invalid k value")
	ErrEmptyArray     = errors.New("array is empty")
	ErrKTooLarge      = errors.New("k is too large")
	ErrKthNotFound    = errors.New("k-th element not found")
)

func readInput() ([]int, int, error) {
	var count int

	_, err := fmt.Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read count: %w", err)
	}

	arr := make([]int, count)

	for index := range count {
		_, err := fmt.Scan(&arr[index])
		if err != nil {
			return nil, 0, fmt.Errorf("failed to read element: %w", err)
		}
	}

	var kValue int

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read k value: %w", err)
	}

	if kValue < 1 || kValue > count {
		return nil, 0, ErrInvalidK
	}

	return arr, kValue, nil
}

func findKthLargest(arr []int, kValue int) (int, error) {
	if len(arr) == 0 {
		return 0, ErrEmptyArray
	}

	if kValue > len(arr) {
		return 0, ErrKTooLarge
	}

	minHeap := InitMinHeap(kValue)
	for _, num := range arr {
		if minHeap.Len() < kValue {
			heap.Push(minHeap, num)
		} else {
			peekValue, err := minHeap.Peek()
			if err != nil {
				return 0, err
			}

			if num > peekValue {
				heap.Pop(minHeap)
				heap.Push(minHeap, num)
			}
		}
	}

	if minHeap.Len() < kValue {
		return 0, ErrKthNotFound
	}

	result, err := minHeap.Peek()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func main() {
	arr, kValue, err := readInput()
	if err != nil {
		fmt.Printf("input error: %v\n", err)

		return
	}

	result, err := findKthLargest(arr, kValue)
	if err != nil {
		fmt.Printf("processing error: %v\n", err)

		return
	}

	fmt.Println(result)
}
