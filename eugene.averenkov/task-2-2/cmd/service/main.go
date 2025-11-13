package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrInvalidN       = errors.New("error invalid number")
	ErrInvalidElement = errors.New("error invalid element")
	ErrInvalidK       = errors.New("error invalid range k")
	ErrEmptyArray     = errors.New("error empty array")
	ErrKTooLarge      = errors.New("error large k")
	ErrKthNotFound    = errors.New("error not found k")
)

func readInput() ([]int, int, error) {
	var count int

	_, err := fmt.Scan(&count)
	if err != nil {
		return nil, 0, ErrInvalidN
	}

	arr := make([]int, count)

	for index := range count {
		_, err := fmt.Scan(&arr[index])
		if err != nil {
			return nil, 0, ErrInvalidElement
		}
	}

	var kValue int

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return nil, 0, ErrInvalidK
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
		} else if num > minHeap.Peek() {
			heap.Pop(minHeap)
			heap.Push(minHeap, num)
		}
	}

	if minHeap.Len() < kValue {
		return 0, ErrKthNotFound
	}

	return minHeap.Peek(), nil
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
