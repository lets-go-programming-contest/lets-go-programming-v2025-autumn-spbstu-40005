package main

import (
	"fmt"
	"os"
)

func readInput() ([]int, int, error) {
	var N int
	_, err := fmt.Scan(&N)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка чтения N: %v", err)
	}

	if N < 1 || N > 10000 {
		return nil, 0, fmt.Errorf("N должно быть в диапазоне от 1 до 10000, получено: %d", N)
	}

	arr := make([]int, N)
	for i := 0; i < N; i++ {
		_, err := fmt.Scan(&arr[i])
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка чтения элемента %d: %v", i+1, err)
		}
		if arr[i] < -10000 || arr[i] > 10000 {
			return nil, 0, fmt.Errorf("элемент %d выходит за диапазон [-10000, 10000]: %d", i+1, arr[i])
		}
	}

	var k int
	_, err = fmt.Scan(&k)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка чтения k: %v", err)
	}

	if k < 1 || k > N {
		return nil, 0, fmt.Errorf("k должно быть в диапазоне от 1 до %d, получено: %d", N, k)
	}

	return arr, k, nil
}

func findKthLargest(arr []int, k int) (int, error) {
	if len(arr) == 0 {
		return 0, fmt.Errorf("пустой массив")
	}
	if k > len(arr) {
		return 0, fmt.Errorf("k (%d) больше длины массива (%d)", k, len(arr))
	}

	h := InitMinHeap()

	for _, num := range arr {
		if h.Len() < k {
			h.PushValue(num)
		} else if num > h.Peek() {
			h.PopValue()
			h.PushValue(num)
		}
	}

	if h.Len() < k {
		return 0, fmt.Errorf("не удалось найти %d-й наибольший элемент", k)
	}

	return h.Peek(), nil
}

func main() {
	arr, k, err := readInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка ввода: %v\n", err)
		os.Exit(1)
	}

	result, err := findKthLargest(arr, k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка обработки: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
