package main

import (
	"container/heap"
	"errors"
	"fmt"

	"nikita.brevnov/task-2-2/internal/intheap"
)

var (
	errEmptyHeap = errors.New("the heap is already empty")
	errConvert   = errors.New("not integer in heap")
)

func findLargest(nums []int, number int) (int, error) {
	dishesHeap := &intheap.IntHeap{}
	heap.Init(dishesHeap)

	for _, num := range nums {
		heap.Push(dishesHeap, num)
	}

	for range number - 1 {
		heap.Pop(dishesHeap)
	}

	val := heap.Pop(dishesHeap)
	if val == nil {
		return 0, errEmptyHeap
	}

	value, ok := val.(int)
	if !ok {
		return 0, errConvert
	}

	return value, nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover, %v\n", r)

			return
		}
	}()

	var countOfDishes int
	if _, err := fmt.Scan(&countOfDishes); err != nil {
		fmt.Printf("invalid input for countOfDishes: %v\n", err)

		return
	}

	nums := make([]int, countOfDishes)
	for index := range countOfDishes {
		if _, err := fmt.Scan(&nums[index]); err != nil {
			fmt.Printf("invalid input for nums: %v\n", err)

			return
		}
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		fmt.Printf("invalid input for number: %v\n", err)

		return
	}

	if number < 1 || number > countOfDishes {
		fmt.Println("this number is not suitable")

		return
	}

	res, err := findLargest(nums, number)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
