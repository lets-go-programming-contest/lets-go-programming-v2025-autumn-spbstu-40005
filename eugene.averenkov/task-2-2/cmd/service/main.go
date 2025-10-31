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

const (
	EdgeArray = 10000
	MinCount  = 1
)

func readInput() ([]int, int, error) {
	var count int

	_, err := fmt.Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка чтения N: %w", err)
	}

	if count < MinCount || count > EdgeArray {
		return nil, 0, fmt.Errorf("%w, получено: %d", ErrInvalidN, count)
	}

	arr := make([]int, count)

	for index := range count {
		_, err := fmt.Scan(&arr[index])
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка чтения элемента %d: %w", index+1, err)
		}

		if arr[index] < -EdgeArray || arr[index] > EdgeArray {
			return nil, 0, fmt.Errorf("%w: элемент %d имеет значение %d", ErrInvalidElement, index+1, arr[index])
		}
	}

	var kValue int

	_, err = fmt.Scan(&kValue)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка чтения k: %w", err)
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
		return 0, fmt.Errorf("%w: k=%d, длина=%d", ErrKTooLarge, kValue, len(arr))
	}

	minHeap := InitMinHeap()

	for _, num := range arr {
		if minHeap.Len() < kValue {
			minHeap.PushValue(num)
		} else if num > minHeap.Peek() {
			minHeap.PopValue()
			minHeap.PushValue(num)
		}
	}

	if minHeap.Len() < kValue {
		return 0, fmt.Errorf("%w: k=%d", ErrKthNotFound, kValue)
	}

	return minHeap.Peek(), nil
}

func main() int {
	arr, kValue, err := readInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка ввода: %v\n", err)
		return 1
	}

	result, err := findKthLargest(arr, kValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка обработки: %v\n", err)
		return 1
	}

	fmt.Println(result)
}
