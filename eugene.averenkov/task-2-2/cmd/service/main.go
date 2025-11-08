package main

import (
	"errors"
	"fmt"
	"os"
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
		return nil, 0, fmt.Errorf("read error N: %w", err)
	}

	arr := make([]int, count)

	for index := range count {
		_, err := fmt.Scan(&arr[index])
		if err != nil {
			return nil, 0, fmt.Errorf("read element error %d: %w", index+1, err)
		}
	}

	var kValue int

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return nil, 0, fmt.Errorf("read error k: %w", err)
	}

	if kValue < 1 || kValue > count {
		return nil, 0, fmt.Errorf("%w: k=%d, N=%d", ErrInvalidK, kValue, count)
	}

	return arr, kValue, nil
}

func findKthLargest(arr []int, kValue int) (int, error) {
	if len(arr) == 0 {
		return 0, ErrEmptyArray
	}

	if kValue > len(arr) {
		return 0, fmt.Errorf("%w: k=%d, length=%d", ErrKTooLarge, kValue, len(arr))
	}

	minHeap := InitMinHeap()

	for _, num := range arr {
		if minHeap.Len() < kValue {
			heap.Push(minHeap, num)
		} else if num > minHeap.Peek() {
			heap.Pop(minHeap)
			heap.Push(minHeap, num)
		}
	}

	if minHeap.Len() < kValue {
		return 0, fmt.Errorf("%w: k=%d", ErrKthNotFound, kValue)
	}

	return minHeap.Peek(), nil
}

func main() {
	arr, kValue, err := readInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %v\n", err)

		return
	}

	result, err := findKthLargest(arr, kValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "processing error: %v\n", err)

		return
	}

	fmt.Println(result)
}
