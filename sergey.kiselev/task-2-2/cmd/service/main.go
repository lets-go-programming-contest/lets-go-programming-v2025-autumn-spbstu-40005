package main

import (
	"container/heap"
	"errors"
	"fmt"

	"sergey.kiselev/task-2-2/internal/intheap"
)

var (
	errOutofRange = errors.New("the heap is already empty")
	errConvert    = errors.New("not integer in heap")
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
		return 0, errOutofRange
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

	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Printf("invalid input for dishesCount: %v\n", err)

		return
	}

	nums := make([]int, dishesCount)
	for index := range dishesCount {
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

	if number > dishesCount {
		fmt.Println("this number is not suitable")

		return
	}

	res, err := findLargest(nums, number)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
